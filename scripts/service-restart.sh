sudo service pass_tsel_server restart

while IFS='=' read -r key value; do
  if [ "$key" = "MO_THREAD" ]; then
    n=$value

    i=1
    while [ $i -le $n ]
    do
      sudo service pass_tsel_consumer_mo@"thread_$i" restart
      sleep 1
      i=$((i+1))
    done
  fi
done < "../.env"

while IFS='=' read -r key value; do
  if [ "$key" = "DR_THREAD" ]; then
    n=$value

    i=1
    while [ $i -le $n ]
    do
      sudo service pass_tsel_consumer_dr@"thread_$i" restart
      sleep 1
      i=$((i+1))
    done
  fi
done < "../.env"