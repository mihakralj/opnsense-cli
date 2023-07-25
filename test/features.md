opnsense
    status
    config

    reboot
    service
    config
        show
        comment
        set
        delete
        save
        load
        sign
        validate
        commit
        rollback
    service



ifconfig -a | sed -E 's/^([a-zA-Z0-9]*:)(.*)/\1\n       \2/; s/=/\: /g' | awk '{ if (NF > 0 && substr($1,length($1)) != ":") {$1 = $1 ":"; print "        " $0} else {print $0}}'


for file in /usr/local/opnsense/service/conf/actions.d/actions_*.conf; do service=$(basename $file | awk -F'[_.]' '{print $2}'); echo "${service}:"; awk -F'[][]' '/^\[/{print "  " $2 ":"; next} {if($1 ~ /:/ && $1 !~ /: /){gsub(/:/,": ",$1)}; if(NF>0) print "    " $1}' $file; done
