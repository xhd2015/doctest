# Scenario

**Feature**: signature rules accept both without-d and with-d shapes for Setup/Run/Assert

```
# parse + rules.Check*
author Setup/Run/Assert with optional second param d *session.Doctest after t
  -> accepted (P2)
author without d
  -> still accepted
```

## Preconditions

- Leaves set `req.Op` to `parse-with-d` or `parse-without-d`.

## Steps

1. Descendant selects parse mode.

## Context

- Today `rules.Check*` only accept without-d → `with-d-accepted` is RED until implementer.
