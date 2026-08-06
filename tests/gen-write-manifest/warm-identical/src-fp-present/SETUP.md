# Scenario

**Feature**: warm WriteGoMod leaves source-input fingerprint on disk

```
first WriteGoMod -> second identical WriteGoMod
  -> doctest.gomod-src present
```

## Steps

1. Inherit warm-identical parent Setup.
2. Assert fingerprint file after measured second write.
