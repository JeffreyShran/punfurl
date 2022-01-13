# punfurl

Takes any url and creates a powerset from it which is useful for fuzzing.

## Example
```
root@workstation:~/tools/punfurl# echo "https://support.google.com/google-ads/answer/2472708?hl=en-GB"  | go run ~/tools/punfurl/main.go
https://support.google.com/google-ads
https://support.google.com/answer
https://support.google.com/google-ads/answer
https://support.google.com/2472708
https://support.google.com/google-ads/2472708
https://support.google.com/answer/2472708
https://support.google.com/google-ads/answer/2472708
https://support.google.com
https://support.google.com/google-ads/answer/2472708?hl=en-GB
```