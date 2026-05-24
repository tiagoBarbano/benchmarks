wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"

local payload = string.rep("abcdefghij", 15000)

wrk.body = [[
{
  "id": 1,
  "name": "benchmark",
  "payload": "]] .. payload .. [["
}
]]