stinkbait
=========

![lunch](./stinkbait.jpg)

Marketing site url checker

How Do
------

Firstly, get yourself a toke (valid for 5 minutes):
```bash
$ curl -XPOST https://stinkbait.in.opsee.com/token
# => {"token":"69a66ab4b6cad36d88159a1b7040f1f27968c3a63da2d0bd5245ea96f6aef087"}
```

Then, use that toke as a bearer token in an authorization header to make a test check request (response bodies will be capped at 4096 bytes):
```bash
curl https://stinkbait.in.opsee.com/check -H 'Authorization: Bearer 69a66ab4b6cad36d88159a1b7040f1f27968c3a63da2d0bd5245ea96f6aef087' -d '{
  "max_hosts": 1,
  "check": {
    "id": "www.reddit.com",
    "interval": 0,
    "target": {
      "name": "www.reddit.com",
      "type": "url",
      "id": "www.reddit.com",
      "address": "www.reddit.com"
    },
    "name": "www.reddit.com",
    "http_check": {
      "name": "",
      "path": "/r/pepe",
      "protocol": "https",
      "port": 443,
      "verb": "GET",
      "body": ""
    }
  }
}'

# => {
  "responses": [
    {
      "target": {
        "name": "www.reddit.com",
        "type": "url",
        "id": "www.reddit.com",
        "address": "www.reddit.com"
      },
      "error": "",
      "passing": false,
      "http_response": {
        "code": 200,
        "body": "<!doctype html><html xmlns=\"http://www.w3.org/1999/xhtml\" lang=\"en\" xml:lang=\"en\"><head><title>Rare Pepes</title><meta name=\"keywords\" content=\" reddit, reddit.com, vote, comment, submit \" /><meta name=\"description\" content=\"reddit: the front page of the internet\" /><meta name=\"referrer\" content=\"always\"><meta http-equiv=\"Content-Type\" content=\"text/html; charset=UTF-8\" /><link rel=\"alternate\" media=\"only screen and (max-width: 640px)\" href=\"https://m.reddit.com/r/pepe\" /><meta name=\"viewport\" content=\"width=1024\"><meta property=\"og:image\" content=\"https://www.redditstatic.com/icon.png\"><meta property=\"og:site_name\" content=\"reddit\"><meta property=\"og:description\" content=\"reddit: the front page of the internet\"><meta property=\"og:title\" content=\"Rare Pepes • /r/pepe\"><meta property=\"twitter:site\" content=\"reddit\"><meta property=\"twitter:card\" content=\"summary\"><meta property=\"twitter:title\" content=\"Rare Pepes • /r/pepe\"><link rel=icon href=\"//www.redditstatic.com/icon.png\" sizes=\"256x256\" type=\"image/png\" /><link rel=shortcut icon href=\"//www.redditstatic.com/favicon.ico\" type=\"image/x-icon\" /><link rel=apple-touch-icon-precomposed href=\"//www.redditstatic.com/icon-touch.png\" /><link rel=\"alternate\" type=\"application/atom+xml\" title=\"RSS\" href=\"http://www.reddit.com:443/r/pepe/.rss\" /><link rel=\"stylesheet\" type=\"text/css\" href=\"//www.redditstatic.com/reddit.UgiHnyijbD8.css\" media=\"all\"><!--[if gte IE 8]><!--><link rel=\"stylesheet\" href=\"https://b.thumbs.redditmedia.com/zKdrUQPUygrwlt1tnTMqtoj9QLQcpzfDetGD2DYiWmQ.css\" title=\"applied_subreddit_stylesheet\" type=\"text/css\"><!--<![endif]--><!--[if gte IE 9]><!--><script type=\"text/javascript\" src=\"//www.redditstatic.com/reddit-init.en.iHNXlWNBwAo.js\"></script><!--<![endif]--><!--[if lt IE 9]><script type=\"text/javascript\" src=\"//www.redditstatic.com/reddit-init-legacy.en.6-XDiAsNODw.js\"></script><![endif]--><script type=\"text/javascript\" id=\"config\">r.setup({\"ajax_domain\": \"www.reddit.com:443\", \"post_site\": \"pepe\", \"gold\": false, \"poisoning_report_mac\": null, \"requires_eu_cookie_policy\": false, \"nsfw_media_acknowledged\": false, \"stats_domain\": \"https://stats.redditmedia.com\", \"feature_expando_events\": true, \"facebook_app_id\": \"322647334569188\", \"is_sponsor\": false, \"user_id\": false, \"eventtracker_url\": \"//pixel.redditmedia.com/pixel/of_delight.png\", \"poisoning_canary\": null, \"logged\": false, \"over_18\": false, \"uitracker_url\": \"//pixel.redditmedia.com/pixel/of_discovery.png\", \"modhash\": \"\", \"send_logs\": true, \"listing_over_18\": false, \"https_endpoint\": false, \"extension\": null, \"cur_site\": \"t5_2qrfy\", \"event_target\": {\"target_id\": 4607854, \"target_type\": \"listing\", \"target_sort\": \"hot\", \"target_fullname\": \"t5_2qrfy\"}, \"tracker_url\": \"//pixel.redditmedia.com/pixel/of_destiny.png?v=EoUNZsPfCBJaNpI7JiHAFUSB%2BITBt8jwFeuJrqfFid82WUEhJ%2F970XJKa5jAoJ%2BlL%2BJb0rPH%2B6w%3D\", \"expando_preference\": null, \"debug\": false, \"has_subscribed\": false, \"server_time\": 1456936441.0, \"clicktracker_url\": \"//pixel.redditmedia.com/click\", \"pref_no_profanity\": true, \"cur_domain\": \"reddit.com:443\", \"feature_new_report_dialog\": true, \"events_collector_url\": \"https://events.redditmedia.com/v1\", \"https_forced\": true, \"user_in_timeout\": false, \"is_fake\": false, \"renderstyle\": \"html\", \"vote_hash\": \"zok+Gv5GRy2JzzFtT2mPdKyCSgQY16pINTMbG5DhIQ7yLqHpmDVcPHUoOpIqnLqM6p0doNCEhLOFH5eihc1YLUp9v+WjW3qQ8DxQwyzJYtEWOVAKPZYI97J2BXa7AxPAwalBo9wXIQQsJYtE4YL3rnNcWUl+0SJIR2uOV5+YYys=\", \"events_collector_secret\": \"thiiM0ahsieSiech1phithe6chahngoo8sah6aid\\n\", \"adtracker_url\": \"//pixel.redditmedia.com/pixel/of_doom.png\", \"events_collector_key\": \"RedditFrontend1\", \"cache_policy\": \"loggedout_www\", \"store_visits\": false, \"anon_eventtracker_url\": \"//pixel.redditmedia.com/pixel/of_diversity.png\", \"new_window\": false, \"pref_beta\": false, \"eu_cookie_max_attempts\": 3, \"pageInfo\": {\"actionName\": \"hot.GET_listing\", \"verification\": \"adfad553bf3b6c42e703170c536073e9c1f952f0\"}, \"mweb_blacklist_expressions\": [\"^/prefs/?\", \"^/live/?\", \"/message/\", \"/wiki/?\", \"/m/\", \"^/subreddits/create\", \"^/gold\", \"^/advertising\", \"^/promoted",
        "headers": [
          {
            "name": "X-Ua-Compatible",
            "values": [
              "IE=edge"
            ]
          },
          {
            "name": "X-Content-Type-Options",
            "values": [
              "nosniff"
            ]
          },
          {
            "name": "X-Xss-Protection",
            "values": [
              "1; mode=block"
            ]
          },
          {
            "name": "Vary",
            "values": [
              "accept-encoding"
            ]
          },
          {
            "name": "Cache-Control",
            "values": [
              "max-age=0, must-revalidate"
            ]
          },
          {
            "name": "X-Moose",
            "values": [
              "majestic"
            ]
          },
          {
            "name": "Server",
            "values": [
              "cloudflare-nginx"
            ]
          },
          {
            "name": "Content-Type",
            "values": [
              "text/html; charset=UTF-8"
            ]
          },
          {
            "name": "Cf-Ray",
            "values": [
              "27d640f2589e2a4f-SEA"
            ]
          },
          {
            "name": "Date",
            "values": [
              "Wed, 02 Mar 2016 16:34:01 GMT"
            ]
          },
          {
            "name": "Strict-Transport-Security",
            "values": [
              "max-age=15552000; includeSubDomains; preload"
            ]
          },
          {
            "name": "X-Frame-Options",
            "values": [
              "SAMEORIGIN"
            ]
          },
          {
            "name": "Connection",
            "values": [
              "keep-alive"
            ]
          },
          {
            "name": "Set-Cookie",
            "values": [
              "__cfduid=d5fc962add44a52df1a392ef61be2028d1456936440; expires=Thu, 02-Mar-17 16:34:00 GMT; path=/; domain=.reddit.com; HttpOnly",
              "loid=fGYTVkYbTvZ8JTXpe6; Domain=reddit.com; Max-Age=63071999; Path=/; expires=Fri, 02-Mar-2018 16:34:01 GMT; secure",
              "loidcreated=2016-03-02T16%3A34%3A01.085Z; Domain=reddit.com; Max-Age=63071999; Path=/; expires=Fri, 02-Mar-2018 16:34:01 GMT; secure"
            ]
          }
        ],
        "metrics": [
          {
            "name": "request_latency_ms",
            "value": 733.061273
          }
        ],
        "host": ""
      }
    }
  ],
  "error": ""
}

```
