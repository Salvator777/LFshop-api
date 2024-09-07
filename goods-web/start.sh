srv_name="goods_web_main"
chmod +x ./$srv_name

# 重启，如果已经存在则关闭重启
if pgrep -x $srv_name > /dev/null
then
  echo "${srv_name} is running"
  echo "shutting down ${srv_name}"
  pkill -x $srv_name
  if [ $? -eq 0 ]; then
    echo "Process $srv_name terminated successfully"
  else
    echo "Failed to terminate process $srv_name"
  fi
  echo "starting ${srv_name}"
  ./$srv_name > /dev/null 2>&1 &
  if [ $? -eq 0 ]; then
    echo "start ${srv_name} success"
  else
    echo "Failed to start ${srv_name}"
  fi
else
  echo "starting ${srv_name}"
  ./$srv_name > /dev/null 2>&1 &
  if [ $? -eq 0 ]; then
    echo "start ${srv_name} success"
  else
    echo "Failed to start ${srv_name}"
  fi
fi
