-- luacheck: globals KEYS ARGV redis
---@diagnostic disable: undefined-global

local KEYS = _G.KEYS
local ARGV = _G.ARGV
local redis = _G.redis

-- Fixed-window rate limit script.
-- KEYS[1]: key
-- ARGV[1]: window_seconds
-- ARGV[2]: max_requests

local key = KEYS[1]
local window = tonumber(ARGV[1])
local max_requests = tonumber(ARGV[2])

local current = redis.call('INCR', key)
if current == 1 then
  redis.call('EXPIRE', key, window)
end

if current > max_requests then
  return 0
end

return 1
