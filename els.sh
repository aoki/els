#!/bin/bash -eu

readonly CACHE_FILE=/tmp/ec2-ls.cache

if [[ $# -ne 0 ]]; then
  case ${1-} in
    -r|--renew)
      rm ${CACHE_FILE}
      ;;
    *|-h|--help)
      printf "List EC2 instances. (If spcify -r option then update cache file.)\nCache file is ${CACHE_FILE}"
      exit 0;;
  esac
fi

if [[ -e ${CACHE_FILE} ]]; then
  cat ${CACHE_FILE}
  printf "\e[31mThis results is using cache. If you want to updating cache, use -r or --renew option.\n\e[m" 1>&2
else
  aws ec2 describe-instances --filter "Name=instance-state-name,Values=running" | \
    jq -r '.Reservations[].Instances[] | (.Tags | select(. != null) | from_entries | select(.Environment != null) | {env: (.Environment), role: (.Role), name: (.Name)}) + {private: (.PrivateIpAddress), public: (.PublicIpAddress), type: (.InstanceType)} |  to_entries| [.[].value] | ., (["Environment", "Role", "Name", "PrivateIP", "PublicIP", "InstanceType"]) | @sh' | \
    sed -e s/\'//g | column -t | sort  | uniq | tee ${CACHE_FILE}
fi
