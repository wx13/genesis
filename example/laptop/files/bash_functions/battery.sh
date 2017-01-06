function battery {
	a=( $(acpi) )
	state=${a[2]}
	percent=${a[3]}
	percent=${percent%,}
	case $state in
	"Discharging,")
		state="B"
		;;
	"Full,")
		state="A"
		;;
	"Charging,")
		state="C"
		;;
	*)
		state=""
	esac
	case $state in
	[BC])
		t=${a[4]}
		t=( ${t//:/ } )
		h=${t[0]}
		m=${t[1]}
		h=${h#0}
		t=${h}:${m}
		;;
	*)
		t=""
	esac
	echo $state    $percent   $t
}
