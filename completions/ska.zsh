#compdef ska

local ret=1

_arguments -C \
	   '1: :->tpl' \
	   '*: :->args' && ret=0

case $state in
    tpl)
	local templates; templates=()
	local skadir=$HOME/.local/share/ska

	for tpl in $(ls $skadir); do
	    if [ -f $skadir/$tpl/values.toml ]; then
		local desc=$(head -n 1 $skadir/$tpl/values.toml)
	    else
		local desc="$tpl has not values.toml"
	    fi
	    templates+="$tpl:$desc"
	done

	templates+="help:Getting help"
	
	 _describe -t templates 'template' templates && ret=0
    ;;
    args)
	local arguments; arguments=(
	    '-e:Editor command (nvim)'
	    '-o:Output directory (.)'
	    '-t:Templates directory ($HOME/.local/share/ska)'
	)

	 _describe -t arguments 'ska args' arguments && ret=0
    ;;
esac

return ret
