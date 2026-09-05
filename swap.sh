if [ -f /apphub/ads.alakis.com/bin/orbit-server-linux-amd64.new ]; then
  mv -f /apphub/ads.alakis.com/bin/orbit-server-linux-amd64.new /apphub/ads.alakis.com/bin/orbit-server-linux-amd64
  chmod +x /apphub/ads.alakis.com/bin/orbit-server-linux-amd64
  systemctl restart orbit
  sleep 4
  echo "--- restarted ---"
  journalctl -u orbit --no-pager -n 15 | grep -E 'connections|Listening|panic|fatal|error' | head -10
else
  echo "no .new file found"
fi
