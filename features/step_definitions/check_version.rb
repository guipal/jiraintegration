require "net/http"
require "uri"
require "rspec/json_expectations"

When(/^I send "([^"]*)" request to "([^"]*)"$/) do |arg1, arg2|
  uri = URI.parse(arg2)

  http = Net::HTTP.new(uri.host, uri.port)
  request = Net::HTTP::Get.new(uri.request_uri)

  @response = http.request(request) 
end

Then(/^the response code should be (\d+)$/) do |arg1|
  @response.code.should == arg1
end

Then(/^the response should match json:$/) do |string|
  jsonObject = JSON.parse(string)
  expect(@response.body).to include_json(jsonObject)
end
