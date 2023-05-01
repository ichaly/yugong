
#curl 'https://edith.xiaohongshu.com/api/sns/web/v1/user/otherinfo?target_user_id=5b797f9b6bd7380001d511b8' \
#  -H 'cookie: web_session=040069b5511a2a147061d4f17a364b16fb5f6c' \
#  -H 'referer: https://www.xiaohongshu.com/' \
#  -H 'x-t: 1682947892201' \
#  -H 'x-s: 0YsC1iavZ2w6O6M+slkkOiT+OYFp1laB0Y1Csidvs6M3' \
#  -H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.54' \
#  --compressed


curl 'https://edith.xiaohongshu.com/api/sns/web/v1/user_posted?num=30&cursor=&user_id=5b797f9b6bd7380001d511b8' \
  -H 'cookie: web_session=040069b5511a2a147061d4f17a364b16fb5f6c' \
  -H 'referer: https://www.xiaohongshu.com/' \
  -H 'X-T: 1682947892666' \
  -H 'X-S: sl5C1iwkOidkZB1b0gwvOgvLOgci125W1l4UslVJ12F3' \
  -H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.54' \
  --compressed

#curl 'https://edith.xiaohongshu.com/api/sns/web/v1/feed' \
#  -X POST -d '{"source_note_id":"63ff46b1000000001300b1b9"}' \
#  -H "Content-Type: application/json" \
#  -H 'cookie: web_session=040069b5511a2a147061d4f17a364b16fb5f6c' \
#  -H 'referer: https://www.xiaohongshu.com/' \
#  -H 'x-t: 1682954648469' \
#  -H 'x-s: O6OvOiMl1lO602Mp1BTi1gMi12T+ZYsl1BkBslsLZY13' \
#  -H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.54' \
#  --compressed