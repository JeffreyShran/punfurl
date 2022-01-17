# punfurl

Takes URL's on `stdin` and jumbles the paths using a powerset technique and recombines with the domain on `stdout`.

## Example
```
~/> echo "https://jeff.com/api/v1/datasources/iModels/8d73d54f/extraction/run" | go run ~/tools/punfurl/main.go                                                             https://jeff.com/api
https://jeff.com/v1
https://jeff.com/api/v1
https://jeff.com/datasources
https://jeff.com/api/datasources
https://jeff.com/v1/datasources
https://jeff.com/api/v1/datasources
https://jeff.com/iModels
https://jeff.com/api/iModels
https://jeff.com/v1/iModels
https://jeff.com/api/v1/iModels
https://jeff.com/datasources/iModels
https://jeff.com/api/datasources/iModels
https://jeff.com/v1/datasources/iModels
https://jeff.com/api/v1/datasources/iModels
https://jeff.com/8d73d54f
https://jeff.com/api/8d73d54f
https://jeff.com/v1/8d73d54f
https://jeff.com/api/v1/8d73d54f
https://jeff.com/datasources/8d73d54f
https://jeff.com/api/datasources/8d73d54f
https://jeff.com/v1/datasources/8d73d54f
https://jeff.com/api/v1/datasources/8d73d54f
https://jeff.com/iModels/8d73d54f
https://jeff.com/api/iModels/8d73d54f
https://jeff.com/v1/iModels/8d73d54f
https://jeff.com/api/v1/iModels/8d73d54f
https://jeff.com/datasources/iModels/8d73d54f
https://jeff.com/api/datasources/iModels/8d73d54f
https://jeff.com/v1/datasources/iModels/8d73d54f
https://jeff.com/api/v1/datasources/iModels/8d73d54f
https://jeff.com/extraction
https://jeff.com/api/extraction
https://jeff.com/v1/extraction
https://jeff.com/api/v1/extraction
https://jeff.com/datasources/extraction
https://jeff.com/api/datasources/extraction
https://jeff.com/v1/datasources/extraction
https://jeff.com/api/v1/datasources/extraction
https://jeff.com/iModels/extraction
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

This implemention *does* record the ordering where it can which is useful for api testing. i.e. We often see `company.xyz/api/v2/doctor/ward` but rarely see `company.xyz/doctor/ward/api/v2` so when we fuzz for `v2` at the end it's a wasted request (or thousand!)

The benefit of this approach is mainly time saving & reduced noise on the target server. Also, adding the logical ordering means that the few tests we complete (versus recursive brute-forcing) have a higher success rate than if it were randomised, although this is more notable for longer URL's than shorter one.

It's not meant to be thorough, it's intended use is for time saving when mass scanning and to be suitable to be pipelined with other open source tooling. On a 6 part path you'd see a 91% reduction of generated URL's using my tool versus bruteforcing with something like FFUF. (64 vs 720)
