package middleware

import (
	"sync"
	"time"

	"werk-ticketing/internal/constants"
	"werk-ticketing/internal/response"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     int           // requests per minute
	burst    int           // burst size
	ttl      time.Duration // time to live for visitor records
}

type visitor struct {
	lastSeen time.Time
	count    int
	mu       sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		burst:    burst,
		ttl:      time.Minute,
	}

	// Cleanup old visitors periodically
	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			v.mu.Lock()
			if time.Since(v.lastSeen) > rl.ttl {
				delete(rl.visitors, ip)
			}
			v.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	v, exists := rl.visitors[ip]
	if !exists {
		v = &visitor{
			lastSeen: time.Now(),
			count:    1,
		}
		rl.visitors[ip] = v
		rl.mu.Unlock()
		return true
	}
	rl.mu.Unlock()

	v.mu.Lock()
	defer v.mu.Unlock()

	// Reset count if a minute has passed
	if time.Since(v.lastSeen) > time.Minute {
		v.count = 1
		v.lastSeen = time.Now()
		return true
	}

	// Check if within rate limit
	if v.count >= rl.rate {
		return false
	}

	v.count++
	v.lastSeen = time.Now()
	return true
}

var globalRateLimiter = NewRateLimiter(constants.RateLimitRequestsPerMinute, constants.RateLimitBurst)
var articleRateLimiter = NewRateLimiter(constants.ArticleRateLimitRequestsPerMinute, constants.ArticleRateLimitBurst)

// tenantRateLimiters stores per-tenant rate limiters
var tenantRateLimiters = make(map[string]*RateLimiter)
var tenantRateLimitersMu sync.RWMutex

// TenantRateLimitConfig defines rate limiting configuration per tenant
type TenantRateLimitConfig struct {
	RequestsPerMinute int
	Burst             int
}

// Default tenant rate limits (can be overridden per tenant)
var defaultTenantRateLimit = TenantRateLimitConfig{
	RequestsPerMinute: 100,
	Burst:             20,
}

// tenantRateLimitConfigs stores custom rate limits per tenant
var tenantRateLimitConfigs = make(map[string]TenantRateLimitConfig)
var tenantRateLimitConfigsMu sync.RWMutex

// SetTenantRateLimit sets custom rate limit for a specific tenant
func SetTenantRateLimit(tenantID string, config TenantRateLimitConfig) {
	tenantRateLimitConfigsMu.Lock()
	defer tenantRateLimitConfigsMu.Unlock()
	tenantRateLimitConfigs[tenantID] = config

	// Reset existing rate limiter so new config takes effect
	tenantRateLimitersMu.Lock()
	delete(tenantRateLimiters, tenantID)
	tenantRateLimitersMu.Unlock()
}

// getTenantRateLimiter returns or creates a rate limiter for a tenant
func getTenantRateLimiter(tenantID string) *RateLimiter {
	tenantRateLimitersMu.RLock()
	rl, exists := tenantRateLimiters[tenantID]
	tenantRateLimitersMu.RUnlock()

	if exists {
		return rl
	}

	// Get tenant-specific config or use defaults
	tenantRateLimitConfigsMu.RLock()
	config, hasCustom := tenantRateLimitConfigs[tenantID]
	tenantRateLimitConfigsMu.RUnlock()

	if !hasCustom {
		config = defaultTenantRateLimit
	}

	// Create new rate limiter for this tenant
	tenantRateLimitersMu.Lock()
	defer tenantRateLimitersMu.Unlock()

	// Double-check after acquiring write lock
	if rl, exists = tenantRateLimiters[tenantID]; exists {
		return rl
	}

	rl = NewRateLimiter(config.RequestsPerMinute, config.Burst)
	tenantRateLimiters[tenantID] = rl
	return rl
}

// RateLimit returns a middleware that rate limits requests
// Skips rate limiting for article endpoints (they have their own rate limiter)
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip rate limiting for article endpoints (they have their own higher limit)
		if c.Request.URL.Path == "/api/v1/articles" || c.Request.URL.Path == "/api/v1/articles/" {
			c.Next()
			return
		}

		ip := c.ClientIP()

		if !globalRateLimiter.allow(ip) {
			response.Error(c, 429, "too many requests")
			c.Abort()
			return
		}

		c.Next()
	}
}

// TenantRateLimit returns a middleware that rate limits requests per tenant
// Each tenant has its own rate limit bucket
func TenantRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get tenant ID from context (set by WithTenant middleware)
		tenantID := GetTenantID(c)
		if tenantID == "" {
			// No tenant, fall back to global rate limiting by IP
			ip := c.ClientIP()
			if !globalRateLimiter.allow(ip) {
				response.Error(c, 429, "too many requests")
				c.Abort()
				return
			}
			c.Next()
			return
		}

		// Use tenant + IP as key for more granular rate limiting
		ip := c.ClientIP()
		key := tenantID + ":" + ip

		rl := getTenantRateLimiter(tenantID)
		if !rl.allow(key) {
			response.Error(c, 429, "too many requests for this tenant")
			c.Abort()
			return
		}

		c.Next()
	}
}

// ArticleRateLimit returns a middleware with higher rate limit for article endpoints
func ArticleRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !articleRateLimiter.allow(ip) {
			response.Error(c, 429, "too many requests")
			c.Abort()
			return
		}

		c.Next()
	}
}
