# sketches-go 
[![Go Reference](https://pkg.go.dev/badge/github.com/graphmetrics/sketches-go.svg)](https://pkg.go.dev/github.com/graphmetrics/sketches-go)
[![Go Version](https://img.shields.io/github/go-mod/go-version/graphmetrics/sketches-go)](https://github.com/graphmetrics/sketches-go)
[![Go Report](https://goreportcard.com/badge/github.com/GraphMetrics/sketches-go)](https://goreportcard.com/report/github.com/graphmetrics/sketches-go)


This repo contains a Go implementation of the distributed quantile sketch algorithm
DDSketch[1] originally developed by [DataDog](https://github.com/datadog/sketches-go)™. DDSketch has relative-error guarantees for any quantile q in [0, 1].
That is if the true value of the qth-quantile is `x` then DDSketch returns a value `y` 
such that `|x-y| / x < e` where `e` is the relative error parameter. DDSketch is also 
fully mergeable, meaning that multiple sketches from distributed systems can be combined 
in a central node.

The default implementation, returned from `NewDefaultDDSketch(relativeAccuracy)`, is
guaranteed not to grow too large in size for any data that can be described by a
distribution whose tails are sub-exponential.

Others implementations are also provided, returned by `LogCollapsingLowestDenseDDSketch(relativeAccuracy, maxNumBins)`
and `LogCollapsingHighestDenseDDSketch(relativeAccuracy, maxNumBins)`, where the q-quantile
will be accurate up to the specified relative error for q that is not too small (or large).
Concretely, the q-quantile will be accurate up to the specified relative error as long as it
belongs to one of the `m` bins kept by the sketch. For instance, If the values are time in seconds, 
`maxNumBins = 2048` covers a time range from 80 microseconds to 1 year.

## References

[1] Charles Masson and Jee E Rim and Homin K. Lee. DDSketch: A fast and fully-mergeable quantile sketch with 
relative-error guarantees. PVLDB, 12(12): 2195-2205, 2019.
