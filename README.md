# uptime-check
Our Pingdom monitor for generating + downloading an image with [HTML/CSS to Image](https://hcti.io).

We have this run every minute and alert us if a check fails.

## Pingdom
Returns results in XML for Pingdom.

```xml
<pingdom_http_custom_check>
  <status>OK</status>
  <response_time>3131</response_time>
</pingdom_http_custom_check>
```

## Deployment
Deploy with `now`.

```
now -e API_KEY=@api_key -e API_ID=@api_id -e UPTIME_PASSWORD=@uptime_password
```
