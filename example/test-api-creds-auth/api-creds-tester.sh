#!/usr/bin/env bash
## requires: python, pyotp, gawk, jq, curl, tofu/terraform


##### arg1 (optional) is the TF Binary (tofu, terraform, /path/to/...) 


## TODO:
##   remove PROXMOX_VE_OTP support
##   test credential presidence in 'fake' section
##   remove less-important fake tests
##   handle python requirements.txt (pyotp)


this_shell_config_file="cred-tester.config.sh"
TF_PLAN_EXTRA_ARGS=()


##### create cred-tester config-shell-file if it does not exist

if [[ -e "${this_shell_config_file}" ]]; then
    source "${this_shell_config_file}"
else
    printf '\ncred tester config file does not exist; creating sample file: "%s"\n'  "${this_shell_config_file}"
    cat <<'EOF' > "${this_shell_config_file}"
#/usr/bin/env bash


##### this file is sourced by the cred-tester script


export PROXMOX_VE_ENDPOINT="${PROXMOX_VE_ENDPOINT:-https://127.0.0.1:8006/}"  ## if unset, use localhost


TF_Output_colorized=1 ## 0 disables color, 1 enables


#### real/fake start/end  overrides are inclusive (allows for skipping tests)
###    test-sets begin numbering with '1'
###    skip real/fake by setting the start-override to a large number

# real_start_override=1
# real_end_override=3

# fake_start_override=1
# fake_end_override=10


curl_connect_timeout=5
# debug_print_env_var_values=0  ## 1 enables; any other value disables
# debug_print_arg_var_values=0  ## 1 enables; any other value disables



#### required: admin creds should have power to unlock a users totp
###    topt user(s) will get locked out due to too many failures
admin_un='admin@pam'
admin_pw='secrets'
admin_totp_s=''



#### test cred sets


###  real creds


real_type_1='api_token [pve]'
real_api_token_1='test1@pve!tokenName=hhhhhhhh-hhhh-hhhh-hhhh-hhhhhhhhhhhh'


real_type_2='user [pam] + pass'
real_un_2='test2@pam'
real_pw_2='secrets'


real_type_3='user [pve] + pass + otp (aka totpS via totp-secret to otp)'
real_un_3='test3@pve'
real_pw_3='secrets'
real_totp_s_3='HHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHH'


# real_type_4='pre-auth auth-ticket + csrf-token'
# real_auth_ticket_4=''
# real_csrf_token_4=''






fake_totp_s='CCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC'

fake_csrf='19870727:TmV2ckdvbm5hR2l2VVVwTmV2ckdvbm5hTGV0VURvd24'
fake_ticket='PVE:base64@pve:19870727::UmljayBBc3RsZXk6IEkganVzdCB3YW5uYSB0ZWxsIHlvdSBob3cgSSdtIGZlZWxpbmcuIEdvdHRhIG1ha2UgeW91IHVuZGVyc3RhbmQuIE5ldmVyIGdvbm5hIGdpdmUgeW91IHVwLiBOZXZlciBnb25uYSBsZXQgeW91IGRvd24uIE5ldmVyIGdvbm5hIHJ1biBhcm91bmQgYW5kIGRlc2VydCB5b3UuIE5ldmVyIGdvbm5hIG1ha2UgeW91IGNyeS4gTmV2ZXIgZ29ubmEgc2F5IGdvb2RieWUuIE5ldmVyIGdvbm5hIHRlbGwgYSBsaWUgYW5kIGh1cnQgeW91Lg=='





### fake creds

## TODO
##   filter out less-interesting tests
##   add all cred methods to tests ( user+pass, api-token, auth-ticket+csrf ) with combinations to fail with [to test credential precidence]


fake_type_1='api_token [fake]'
fake_api_token_1='fake@pve!faker=aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee'


fake_type_2='un+pw [real] + totpS [real] + csrf [fake]'
fake_un_2=$real_un_3
fake_pw_2=$real_pw_3
fake_totp_s_2=$real_totp_s_3
fake_csrf_2=$fake_csrf


fake_type_3='ticket [fake] + csrf [fake]'
fake_auth_ticket_3=$fake_ticket
fake_csrf_token_3=$fake_csrf


fake_type_4='user [real pam] + pass [fake]'
fake_un_4=$real_un_2
fake_pw_4='this_dont_work'


fake_type_5='user [fake pam] + pass [fake]'
fake_un_5='nope@pam'
fake_pw_5='this_dont_work'


fake_type_6='user [fake pve] + pass [fake]'
fake_un_6='nope@pve'
fake_pw_6='this_dont_work'


## less interesting-ish

fake_type_7='user [pve] + pass + otp [fake] (aka totpS via totp-secret to otp)'
fake_un_7=$real_un_3
fake_pw_7=$real_pw_3
fake_totp_s_7=$fake_totp_s


fake_type_8='user [pve] + pass - otp (aka totpS via totp-secret to otp)'
fake_un_8=$real_un_3
fake_pw_8=$real_pw_3


fake_type_9='user [pve] + pass [fake] + otp [fake] (aka totpS via totp-secret to otp)'
fake_un_9=$real_un_3
fake_pw_9='this_dont_work'
fake_totp_s_9=$fake_totp_s


fake_type_10='user [fake pve] + pass [fake] + otp [fake] (aka totpS via totp-secret to otp)'
fake_un_10='nope@pve'
fake_pw_10='this_dont_work'
fake_totp_s_10=$fake_totp_s

EOF

    
    exit
fi






##### helper funcs



### same logic as makefile to determine tofu or terraform
### then override with arg1 if applicable
if [[  "$(tofu -version 2>/dev/null)" == "" ]]; then
    TF_APP=terraform
else
    TF_APP=tofu
fi
export TF_APP="${1:-"${TF_APP}"}"
## TODO: check app is run-able (-x needs full path)
#if [[ ! -x "${TF_APP}" ]]; then
#    printf 'ERROR: TF App is not executable or does not exit: %s\n'  "${TF_APP}"
#    exit 1
#fi





### function to print TF arg-var names for in-line provider testing
###   optional toggle-var (debug_print_arg_var_values) to print values
print_proxmox_arg_var_names() {
    ## required argvars/list of args sent to TF ( "${@}" )
    printf '    set arg-vars:  '
    if [[ -n "${debug_print_arg_var_values}" && "${debug_print_arg_var_values}" == 1 ]]; then
	## 
        printf '%s  '  $(  for x in "$@"; do if [[ "$x" != "-var" ]]; then echo "${x}"; fi ; done | sort  )
    else
        printf '%s  '  $(  for x in "$@"; do if [[ "$x" != "-var" ]]; then echo "${x%%=*}"; fi ; done | sort  )
    fi
    printf '\n'
    
}





### function to print TF variable names (excluding PROXMOX_VE_ENDPOINT)
###   optional toggle-var (debug_print_env_var_values) to print values
print_proxmox_env_var_names() {
    printf '    set cred-vars:  '
    if [[ -n "${debug_print_env_var_values}" && "${debug_print_env_var_values}" == 1 ]]; then
        printf '%s  '  $(env -0 | sort -z | xargs -0 -n1 | awk -F'=' '/^PROXMOX_/ && !/PROXMOX_VE_ENDPOINT/ && $2 > 0  {print $0}')
    else
        printf '%s  '  $(env -0 | sort -z | xargs -0 -n1 | awk -F'=' '/^PROXMOX_/ && !/PROXMOX_VE_ENDPOINT/ && $2 > 0  {print $1}')
    fi
    printf '\n'
}





### function to send the unlock api call to a particular user
### especially handy when testing a legit user and fake pass (too many failures locks the totp)
curl_api_unlock_user_tfa() {
    ## required argvar: username@realm
    local userid="$1" ## user@realm
    local ticket
    local csrf

    if [[ -z "${userid}" ]]; then return; fi

    ## use read to set variables 'ticket' and 'csrf' from the bash-func api_auth with env-vars set and totpx
    read -d $'\0'  ticket  csrf  <<<"$( PROXMOX_VE_USERNAME=${admin_un} PROXMOX_VE_PASSWORD=${admin_pw} api_auth $( totp_secret_to_otp_func "${admin_totp_s}" ) )"
    curl --connect-timeout ${curl_connect_timeout:-5}  -q -s -k -H "X-API-AUTH-TOKEN: ${ticket}"  -H "X-CSRF-TOKEN: ${csrf}" -X PUT "${PROXMOX_VE_ENDPOINT}api2/json/access/users/${userid}/unlock-tfa"
}





### function calling python to generate a TOTP token from a totp-secret
totp_secret_to_otp_func(){  STR=${1} python -c 'import pyotp; import os; import time;  STR=os.getenv("STR", ""); totp = pyotp.TOTP(STR);  print( totp.now() );'  ;  }





### function to call proxmox-api with a user+pass, and optional totp to get an auth_ticket and csrf_token
### endpoint, username, password are passed in via provider env-vars (PROXMOX_VE_ENDPOINT, PROXMOX_VE_USERNAME, PROXMOX_VE_PASSWORD)
api_auth() {
    ## optional arg1: totp value
    local _user_totp_password=$1
    local proxmox_api_ticket_path='api2/json/access/ticket' ## cannot have double "//" - ensure endpoint ends with a "/" and this string does not begin with a "/", or vice-versa

    local api_resp=$( curl -q -s -k --data-urlencode "username=${PROXMOX_VE_USERNAME}"  --data-urlencode "password=${PROXMOX_VE_PASSWORD}"  "${PROXMOX_VE_ENDPOINT}${proxmox_api_ticket_path}" )
    local auth_ticket=$( jq -r '.data.ticket' <<<"${api_resp}" )
    local auth_csrf=$( jq -r '.data.CSRFPreventionToken' <<<"${api_resp}" )

    if [[ $(jq -r '.data.NeedTFA' <<<"${api_resp}") == 1 ]]; then
        api_resp=$( curl -q -s -k  -H "CSRFPreventionToken: ${auth_csrf}" --data-urlencode  "username=${PROXMOX_VE_USERNAME}" --data-urlencode "tfa-challenge=${auth_ticket}" --data-urlencode "password=totp:${_user_totp_password}"  "${PROXMOX_VE_ENDPOINT}${proxmox_api_ticket_path}" )
        auth_ticket=$( jq -r '.data.ticket' <<<"${api_resp}" )
        auth_csrf=$( jq -r '.data.CSRFPreventionToken' <<<"${api_resp}" )
    fi

    printf '%s\n'  "${auth_ticket}"  "${auth_csrf}"
}





### optional color flag for TF
if [[ -n "${TF_Output_colorized}" && "${TF_Output_colorized}" == 0 ]]; then
    TF_PLAN_EXTRA_ARGS+=('-no-color')
fi




### function which calls TF cmd
### arg-order matters; first 3 args are required test-output; remaining args are passed directly to TF
test_cmd() {
    ## required: pass/fail  itr_num  itr_part_num

    local should_pass_or_fail=${1:pass}
    shift
    local test_iteration=${1:-itr-1}
    shift
    local test_part=${1:-part-1}
    shift
    local ret_expt=0
    local res
    local ret
    local timenow=$(date --iso=s)
    timenow=${timenow//:/_}

    res=$( "${TF_APP}" plan "${TF_PLAN_EXTRA_ARGS[@]}"  "$@" 2>&1 )
    ret=$?

    local msg_pass
    local msg_fail

    #msg_pass="passed ${should_pass_or_fail}-as-expected" # successfully pass-as-expected
    #msg_fail="failed ${should_pass_or_fail}-as-expected" # failed pass-as-expected
    if [[ "${should_pass_or_fail,,}" == "pass" ]]; then ## TODO: word these better
        msg_pass="passed ${should_pass_or_fail}-check (expected pass and passed)" # successfully pass-as-expected
        msg_fail="failed ${should_pass_or_fail}-check (expected pass; got an error) " # failed pass-as-expected
    else
        msg_pass="passed ${should_pass_or_fail}-check (expected error and errored)" # successfully fail-as-expected
        msg_fail="failed ${should_pass_or_fail}-check (expected error; got no-errors)" # failed fail-as-expected
    fi

    local got_pass_or_fail

    
    if [[ "${should_pass_or_fail,,}" == "pass"  &&  "${ret}" == 0 ]]  ||  [[ "${should_pass_or_fail,,}" != "pass"  &&  "${ret}" != 0 ]]; then
        printf '  test %i part %i %s\n' ${test_iteration}  ${test_part}  "${msg_pass}"
	got_pass_or_fail='passed'
    else
        printf '  test %i part %i %s\n' ${test_iteration}  ${test_part}  "${msg_fail}"
	got_pass_or_fail='failed'
    fi

    
    local logfile="outs_cred-tester__expect_${should_pass_or_fail}__${got_pass_or_fail}.${timenow}.log"

    printf '%s\n' "${res}"  > "${logfile}" ## write TF output to logfile
    printf '  see TF-output in file: %s\n'  "${logfile}"

    ## use gawk to capture the Errors and print to screen; colored vs not makes a monster
    if [[ -n "${TF_Output_colorized}" && "${TF_Output_colorized}" == 0 ]]; then
	<<<"${res}"  gawk -v out_spaces="  "   '/^Error:/ { S=1; BLOB=out_spaces $0 ; next }   /^$/ {  if(S==1){ nl++; if(nl>1){S=0;nl=0; print BLOB"\n" }else{ BLOB=BLOB"\n"out_spaces $0} }  }   !/^$/ { if(S==1){  BLOB=BLOB"\n"out_spaces $0} }    END { if(S==1){ print BLOB } } '
    else
        ## ansi-color friendly awk'er for tf-outputs (print Error blocks; can be expanded to include any word (eg Warning); change: "+Error:" to  "+([A-Za-z]+:"  )
        ## print only the errors ( use outs_spaces ) as line prefix (expect to be empty when test-for-failure and plan-passes (eg no tf-files)
        <<<"${res}"  gawk -v outs_spaces="  " -v str_bar=$'\xe2\x94\x82' -v str_bar_end=$'\xe2\x95\xb5' '{  raw=$0;  gsub(/\x1b\[([0-9]{1,2}(;[0-9]{1,2})?)?[mGK]/,"");  }  {if(S==1 && match($0,str_bar_end)) {S=0; print outs_spaces raw}  } ;  S==1 {print  outs_spaces raw }       { if(match($0, "^"str_bar" +Error:")) { S=1; print outs_spaces raw }    } '
    fi
        

}





printf 'test endpoint: %s\n'  "${PROXMOX_VE_ENDPOINT}"





#### test real+fake   UN + PW



### message formats + hard-code number of check-types per auth-set
_set_parts=4 ## aka number of "checks" per auth-set
real_test_msg_fmt='\npass-test (expect auth-success): set %i/%i %s/'${_set_parts}';  using %s %s:  %s\n'
fake_test_msg_fmt='\nfail-test (expect auth-failure): set %i/%i %s/'${_set_parts}';  using %s %s:  %s\n'
## expect:  set-#   fake/real_end-#    "part 1/2"   "env OR in-line"   "raw-credentials OR auth-ticket"  ${type}
## todo: ensure comment fmt-opts are correct



test_pre_auth_inline_empty_creds=0  ## want to test this once per run (many 'fake' tests cause this to be re-test if not caught)





#### while loop determining real/fake "end" numbers; starting at 1, checks consecutive vars and bails when both are set
real_start=1
real_end=
fake_start=1
fake_end=


_i=1
while :; do ## infinite loop; maybe break after i gets too high?
    _real="real_type_${_i}"
    _fake="fake_type_${_i}"

    ## bash indirect expansion / dynamic variables
    ## check if _real is a variable and that the value real expands to is not an empty value
    ## once i (real_type_2) is out of bounds/does not exist, sets real_too_far
    ## repeat for fake_type
    ## once both real and fake 'too_far' variables are set, break out of loop
    if [[  -v "${_real}" && -n "${!_real}" ]]; then
	real_end=$_i
    else
	real_too_far=$_i
    fi

    if [[ -v "${_fake}" && -n "${!_fake}" ]]; then
	fake_end=$_i
    else
	fake_too_far=$_i
    fi

    if [[ -n "${real_too_far}" && -n "${fake_too_far}" ]]; then
	break
    fi

    ((_i++))
done
unset _i _real _fake real_too_far fake_too_far








#### nested for-loop
## first for-loop over "real" and "fake" as variable name prefixes
##   then loop over each real/fake credential set
## cred tests order:
##   raw-creds via env-vars
##   raw-creds via in-line (using tf provider-block)
##   pre-auth via env-vars
##   pre-auth via in-line (using tf provider-block)


for loop_real_fake in "real" "fake" ; do

    if [[ "${loop_real_fake}" == "real" ]]; then
	loop_rf_pass_fail="pass"
    else
	loop_rf_pass_fail="fail"
    fi


    loop_rf_test_msg_fmt="${loop_real_fake}_test_msg_fmt" ## eg: loop_rf_test_msg_fmt="real_test_msg_fmt" (literal text)
    loop_rf_test_msg_fmt="${!loop_rf_test_msg_fmt}" ## eg: loop_rf_test_msg_fmt= value of real_test_msg_fmt (variable expansion)
    
    loop_rf_start="${loop_real_fake}_start"
    loop_rf_start="${!loop_rf_start}"

    loop_rf_end="${loop_real_fake}_end"
    loop_rf_end="${!loop_rf_end}"

    loop_rf_start_override="${loop_real_fake}_start_override"
    loop_rf_start_override="${!loop_rf_start_override}"

    loop_rf_end_override="${loop_real_fake}_end_override"
    loop_rf_end_override="${!loop_rf_end_override}"
    

    

    
    printf '\nShould %s:\n' "${loop_rf_pass_fail^}"
    for (( i=${loop_rf_start}; i < ((loop_rf_end+1)) ; i++ )); do
	#### if conditions: allow start/end overrides to skip select tests in order
	if [[ -n "${loop_rf_start_override}" &&  "${loop_rf_start_override}" -gt "${i}" ]]; then
            continue
	fi
	if [[ -n "${loop_rf_end_override}" && "${loop_rf_end_override}" -lt "${i}" ]]; then
            printf '\nscript loop complete via real_end_override; range %i to %i inclusive\n'  ${loop_rf_start_override}  ${loop_rf_end_override}
            break
	fi

	## enforce vars are not set
	unset  PROXMOX_VE_USERNAME  PROXMOX_VE_PASSWORD  PROXMOX_VE_OTP  PROXMOX_VE_API_TOKEN  PROXMOX_VE_AUTH_TICKET  PROXMOX_VE_CSRF_PREVENTION_TOKEN  test_cmd_args


	## bash indirect expansions for fake/real values set in config (aka dynamic variables)
	type="${loop_real_fake}_type_${i}" ## type="real_type_1
	type="${!type}" ## indirect expansion, type set to value from real_type_1
	un="${loop_real_fake}_un_${i}"
	un="${!un}"
	pw="${loop_real_fake}_pw_${i}"
	pw="${!pw}"
	totp_s="${loop_real_fake}_totp_s_${i}"
	totp_s=${!totp_s}
	auth_ticket="${loop_real_fake}_auth_ticket_${i}"
	auth_ticket="${!auth_ticket}"
	csrf_token="${loop_real_fake}_csrf_token_${i}"
	csrf_token="${!csrf_token}"
	apitok="${loop_real_fake}_api_token_${i}"
	apitok="${!apitok}"
	



	printf "\n##"


	### raw-creds via env-vars
	printf "${loop_rf_test_msg_fmt}"  $i  ${loop_rf_end}  "part 1"  "env"  "raw-creds"  "${type}"
	curl_api_unlock_user_tfa "${un}"

    
	if [ -n "${totp_s}" ]; then
            export PROXMOX_VE_OTP=$( totp_secret_to_otp_func "${totp_s}" )
	fi

	export PROXMOX_VE_USERNAME=${un}
	export PROXMOX_VE_PASSWORD=${pw}
	if [[ -n "${apitok}" ]]; then
	    export PROXMOX_VE_API_TOKEN="${apitok}"
	fi
	print_proxmox_env_var_names
	print_proxmox_arg_var_names  "${test_cmd_args[@]}"

	test_cmd ${loop_rf_pass_fail} ${i} 1 ## test with args: "real/fake" "itr" and part "1"





	### raw-creds via in-line
	printf "${loop_rf_test_msg_fmt}"  $i  ${loop_rf_end}  "part 2"  "in-line"  "raw-creds"  "${type}"
	curl_api_unlock_user_tfa "${un}"

	###  array of TF cli credential-arguments (to ultimately set provider-config vars:
	## username, password, otp, api_token
	test_cmd_args=(
	    ## ternary-ish; if cred-type var is set, create the TF cli-arg to pass in a variable override
	    ## eg: if this loop's 'un' is set, append to the arg-array:   -var  username="$un"
	    # $( [[ -n "${auth_ticket}" ]] && printf '%s\n'  "-var"  auth_ticket="${auth_ticket}" ) ## consolidate test types would include ticket+csrf here
	    # $( [[ -n "${csrf_token}" ]] && printf '%s\n'  "-var"  csrf_prevention_token="${csrf_prevention_token}" )
	    $( [[ -n "${un}" ]] && printf '%s\n'  "-var"  username="${un}" )
	    $( [[ -n "${pw}" ]] && printf '%s\n'  "-var"  password="${pw}" )
	    $( [[ -n "${totp_s}" ]] && printf '%s\n'  "-var"  otp="$( totp_secret_to_otp_func "${totp_s}" )" )
	    $( [[ -n "${apitok}" ]] && printf '%s\n'  "-var"  api_token="${apitok}" )
	)

	unset  PROXMOX_VE_USERNAME  PROXMOX_VE_PASSWORD  PROXMOX_VE_API_TOKEN  PROXMOX_VE_OTP
	print_proxmox_env_var_names
	print_proxmox_arg_var_names  "${test_cmd_args[@]}"

	test_cmd ${loop_rf_pass_fail} ${i} 2 "${test_cmd_args[@]}"





	### pre-auth via env-vars
	printf "${loop_rf_test_msg_fmt}"  $i  ${loop_rf_end} "part 3"  "env"  "pre-auth"  "${type}"
	if [[ -n "${apitok}" ]] && [[ -z "${un}" && -z "${PROXMOX_VE_AUTH_TICKET}" ]]; then
	    printf '  skip testing api_token with pre-auth ticket+csrf; invalid test (cannot get auth-ticket and csrf with api-token)\n'
	else
	    curl_api_unlock_user_tfa "${un}"


	    ## call the api_auth function with the totp-password generated from the totp-secret
	    read -d $'\0'  pre_auth_ticket  pre_auth_csrf  <<<"$( PROXMOX_VE_USERNAME="${un}"  PROXMOX_VE_PASSWORD="${pw}"  api_auth $( totp_secret_to_otp_func "${totp_s}" ) )"    
	    unset  PROXMOX_VE_AUTH_TICKET  PROXMOX_VE_CSRF_PREVENTION_TOKEN  PROXMOX_VE_USERNAME  PROXMOX_VE_PASSWORD  PROXMOX_VE_API_TOKEN  PROXMOX_VE_OTP  test_cmd_args
	
	    if [[ -n "${pre_auth_ticket}" && "${pre_auth_ticket}" != "null" ]]; then
		export PROXMOX_VE_AUTH_TICKET="${pre_auth_ticket}"
	    fi
	    if [[ -n "${pre_auth_csrf}" && "${pre_auth_csrf}" != "null" ]]; then
		export PROXMOX_VE_CSRF_PREVENTION_TOKEN="${pre_auth_csrf}"
	    fi
	    if [[ -n "${apitok}" ]]; then
		export PROXMOX_VE_API_TOKEN="${apitok}"
	    fi
	    print_proxmox_env_var_names
	    print_proxmox_arg_var_names  "${test_cmd_args[@]}"
	
	    test_cmd ${loop_rf_pass_fail} ${i} 3
	fi




	### pre-auth via in-line
	printf "${loop_rf_test_msg_fmt}"  $i  ${loop_rf_end} "part 4"  "in-line"  "pre-auth"  "${type}"

	## if apitoken is set and username and env-var auth_ticket are unset, skip test
	## else-if skip test if test-scenario has been tested before (case when expect-to-fail and auth_ticket or csrf are empty/unset
	if [[ -n "${apitok}" ]] && [[ -z "${un}" && -z "${PROXMOX_VE_AUTH_TICKET}" ]]; then
	    printf '  skip testing api_token with pre-auth ticket+csrf; invalid test (cannot get auth-ticket and csrf with api-token)\n'
	elif [[ "${test_pre_auth_inline_empty_creds}" != 0 && "${loop_rf_pass_fail}" == "fail" ]] && [[ -z "${PROXMOX_VE_AUTH_TICKET}" || -z "${PROXMOX_VE_CSRF_PREVENTION_TOKEN}" ]]; then
	    printf '  skip testing: fail test has no pre-auth creds and was already tested with no creds; see std-out for: test_pre_auth_inline_empty_creds\n'
	else
	    printf '  Note: setting test_pre_auth_empty_creds as to not re-test empty pre-auth creds with in-line config\n'
	    test_pre_auth_inline_empty_creds=1 ## flag with var to test no more than once this condition (when pre-auth creds are empty/unset with in-line config)
	    curl_api_unlock_user_tfa "${un}"

	    test_cmd_args=(
		## if auth_ticket / csrf is not set, set to values generated in test above
		$( [[ -z "${auth_ticket}" ]] && printf '%s\n'  "-var"  auth_ticket="${PROXMOX_VE_AUTH_TICKET}" ) ## assume auth-ticket above in REAL was generated successfully
		$( [[ -z "${csrf_token}" ]] && printf '%s\n'  "-var"  csrf_prevention_token="${PROXMOX_VE_CSRF_PREVENTION_TOKEN}" )
		## if auth_ticket / csrf is set, use them
		$( [[ -n "${auth_ticket}" ]] && printf '%s\n'  "-var"  auth_ticket="${auth_ticket}" )
		$( [[ -n "${csrf_token}" ]] && printf '%s\n'  "-var"  csrf_prevention_token="${csrf_prevention_token}" )
	    )

	    unset  PROXMOX_VE_AUTH_TICKET  PROXMOX_VE_CSRF_PREVENTION_TOKEN
	    print_proxmox_env_var_names
	    print_proxmox_arg_var_names  "${test_cmd_args[@]}"

	    test_cmd ${loop_rf_pass_fail} ${i} 4 "${test_cmd_args[@]}"
	fi




	printf '\n'



    done ## end real/fake sets


done ## end "real" vs "fake" variable-prefixes



