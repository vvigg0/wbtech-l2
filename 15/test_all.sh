echo START
pwd
echo hello world
echo $HOME > t.txt && cat < t.txt | wc -c && rm -f t.txt
true && echo AND_OK
false || echo OR_OK
printf "a\nb\n" | grep -q a && echo PIPE_AND_OK || echo PIPE_AND_BAD
printf "a\nb\n" | grep -q z && echo PIPE_OR_BAD || echo PIPE_OR_OK
echo "a | b && c" | wc -c
echo DONE