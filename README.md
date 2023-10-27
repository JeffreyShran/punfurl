# `punfurl`

Takes URL's on `stdin` and jumbles the paths using a powerset technique and recombines with the domain on `stdout`.

## Example
```
~/> echo "https://jeff.com/api/v1/datasources/iModels/8d73d54f/extraction/run" | go run ~/tools/punfurl/main.go
https://jeff.com/api
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

To fuzz directories you typically need to initially create a wordlist then iterate through each word in that list recursively. This takes time and generates lots of requests on the server when you test for them.

`punfurl` essentially take the words in the path (e.g. `api/v1/datasources/iModels/8d73d54f/extraction/run` in the example above) and uses them as the wordlist. However instead of resursively trying every combination `punfurl` uses them to generate a [powerset]("https://en.wikipedia.org/wiki/Power_set") which is a set of all the subsets of a set.

For the set {a,b,c}:

 - The empty set {} is a subset of {a,b,c}
 - And these are subsets: {a}, {b} and {c}
 - And these are also subsets: {a,b}, {a,c} and {b,c}
 - And {a,b,c} is a subset of {a,b,c}
 - And altogether we get the Power Set of {a,b,c}:

P(S) = { {}, {a}, {b}, {c}, {a, b}, {a, c}, {b, c}, {a, b, c} }

Think of it as all the different ways we can select the items (the order of the items doesn't matter), including selecting none, or all.

This implemention *does* record the ordering where it can which is useful for api testing. i.e. We often see `company.xyz/api/v2/doctor/ward` but rarely see `company.xyz/doctor/ward/api/v2` so when we send requests with nonsensical paths it's a wasted request (or thousand!).

The benefit of this approach is mainly to save time & reduce noise on the target server whilst still keeping some logic to the output. The longer the url the more results you'll have.

> If you want to add custom words, just add another *slash* *word* at the end of the URL on `stdin`.
> For example: `company.xyz/api/v2/doctor/ward*/custom_word*`

It's not meant to be thorough, it's intended use is for time saving when mass scanning and to be suitable to be pipelined with other open source tooling such as processing results from tools like GAU and piping into something like httpx.
