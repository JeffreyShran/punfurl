# punfurl

Takes URL's on `stdin` and jumbles the paths using a powerset technique and recombines them to the domain on `stdout`.

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

To fuzz directories you typically need to initially create a worldist then iterate through each word in that list recursively. This takes time and generates lots of requests on the server when you test for them.

punfurl essentially take the words in the path (e.g. `google-ads/answer/2472708` in the example above) and uses them as the wordlist. However, instead of resursively trying every combination, it only takes the [powerset]("https://en.wikipedia.org/wiki/Power_set") which is a set of all the subsets of a set.

For the set {a,b,c}:

 - The empty set {} is a subset of {a,b,c}
 - And these are subsets: {a}, {b} and {c}
 - And these are also subsets: {a,b}, {a,c} and {b,c}
 - And {a,b,c} is a subset of {a,b,c}
 - And altogether we get the Power Set of {a,b,c}:

P(S) = { {}, {a}, {b}, {c}, {a, b}, {a, c}, {b, c}, {a, b, c} }

Think of it as all the different ways we can select the items (the order of the items doesn't matter), including selecting none, or all.

This implemention also records the ordering where it can which is useful for api testing. i.e. We often see `company.xyz/api/v2/doctor/ward` but rarely see `company.xyz/doctor/ward/api/v2` so when we fuzz for `v2` at the end it's a wasted request (or thousand!)

