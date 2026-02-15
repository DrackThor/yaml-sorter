# yaml-sort examples

Before/after examples from simple to advanced, including list-of-objects sorting with a config file.

---

## 1. Simple flat mapping

All keys are sorted alphabetically at the top level.

**Before**

```yaml
zebra: 3
apple: 1
banana: 2
```

**Command:** `yaml-sort file.yaml`

**After**

```yaml
apple: 1
banana: 2
zebra: 3
```

---

## 2. Nested mapping (recursive)

Keys are sorted at every level; nesting is preserved.

**Before**

```yaml
server:
  port: 8080
  host: localhost
  env: prod
database:
  user: admin
  name: mydb
```

**After**

```yaml
database:
  name: mydb
  user: admin
server:
  env: prod
  host: localhost
  port: 8080
```

---

## 3. Mapping with arrays

Sequence (list) elements are not reordered; only mapping keys are sorted.
Keys that hold arrays are sorted by key name like any other.

**Before**

```yaml
tags:
  - a
  - b
  - c
name: myapp
version: "1.0"
```

**After**

```yaml
name: myapp
tags:
  - a
  - b
  - c
version: "1.0"
```

---

## 4. Kubernetes manifest root order (`-k`)

With `-k`, the **root** keys follow a fixed order: `apiVersion`, `kind`, `metadata`, `spec`, `data`, `status`, then the rest alphabetically. Everything under them is still sorted alphabetically (recursive).

**Before**

```yaml
spec:
  replicas: 2
  selector:
    matchLabels:
      app: nginx
metadata:
  name: nginx
  namespace: default
kind: Deployment
apiVersion: apps/v1
```

**Command:** `yaml-sort -k file.yaml`

**After**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
  namespace: default
spec:
  replicas: 2
  selector:
    matchLabels:
      app: nginx
```

---

## 5. List of objects (no config)

Without a config file, a **list of objects** is not reordered: only the **keys inside each object** are sorted alphabetically.
So the list order stays as in the file.

**Before**

```yaml
spec:
  egress:
    - name: nv.egress-1
      action: allow
      ports: any
    - name: nv.egress-0
      action: allow
      ports: any
```

**Command:** `yaml-sort file.yaml` (no `-c`)

**After**

```yaml
spec:
  egress:
    - action: allow
      name: nv.egress-1
      ports: any
    - action: allow
      name: nv.egress-0
      ports: any
```

Keys inside each list item are sorted (`action`, `name`, `ports`), but the list order is unchanged (egress-1 still before egress-0).

---

## 6. List of objects sorted by key (config file `-c`)

To **sort the list itself** by a field (e.g. `name`), use a **config file** and `-c`.
The list at the given path is sorted by that key in each element.

**Config file** (e.g. `.yaml-sort.yaml`):

```yaml
listSortKeys:
  - path: spec.egress
    key: name
```

**Before**

```yaml
spec:
  egress:
    - name: nv.consul-server.consul-egress-1
      action: allow
      ports: tcp/8502
    - name: nv.consul-server.consul-egress-0
      action: allow
      ports: any
```

**Command:** `yaml-sort -c .yaml-sort.yaml file.yaml`

**After**

```yaml
spec:
  egress:
    - action: allow
      name: nv.consul-server.consul-egress-0
      ports: any
    - action: allow
      name: nv.consul-server.consul-egress-1
      ports: tcp/8502
```

The list `spec.egress` is now ordered by `name` (egress-0 before egress-1), and keys inside each item are still alphabetical.

---

## 7. Multiple lists with different sort keys

You can define several rules: each path can use a different key (or the same key).

**Config file** (e.g. `.yaml-sort.yaml`):

```yaml
listSortKeys:
  - path: spec.egress
    key: name
  - path: spec.ingress
    key: name
  - path: spec.process
    key: name
```

**Before**

```yaml
spec:
  ingress:
    - name: nv.ui-ingress-1
      action: allow
    - name: nv.ui-ingress-0
      action: allow
  egress:
    - name: nv.egress-1
      action: allow
    - name: nv.egress-0
      action: allow
  process:
    - name: nginx
      path: /usr/sbin/nginx
    - name: pause
      path: /pause
```

**Command:** `yaml-sort -c .yaml-sort.yaml file.yaml`

**After**

```yaml
spec:
  egress:
    - action: allow
      name: nv.egress-0
    - action: allow
      name: nv.egress-1
  ingress:
    - action: allow
      name: nv.ui-ingress-0
    - action: allow
      name: nv.ui-ingress-1
  process:
    - name: nginx
      path: /usr/sbin/nginx
    - name: pause
      path: /pause
```

- Root keys under `spec` are alphabetical: `egress`, `ingress`, `process`.
- Each of the three lists is sorted by its `name` field.

---

## 8. Kubernetes root order + list sort

Combine `-k` (K8s root order) and `-c` (list sort keys) for manifests that have both a K8s-like root and lists of objects (e.g. NeuVector `NvSecurityRule`).

**Config file** (e.g. `.yaml-sort.yaml`):

```yaml
listSortKeys:
  - path: spec.egress
    key: name
  - path: spec.ingress
    key: name
  - path: spec.process
    key: name
```

**Before**

```yaml
kind: NvSecurityRule
apiVersion: neuvector.com/v1
metadata:
  namespace: publishing-company
  name: nv.bookstore-ui.publishing-company
spec:
  egress:
    - name: nv.consul-server.consul-egress-1
      action: allow
      ports: tcp/8502
    - name: nv.consul-server.consul-egress-0
      action: allow
      applications:
        - Consul
        - SSL
  ingress:
    - name: nv.ui-ingress-1
      action: allow
    - name: nv.ui-ingress-0
      action: allow
  process_profile:
    mode: Protect
    baseline: zero-drift
```

**Command:** `yaml-sort -k -c .yaml-sort.yaml -o sorted.yaml file.yaml`

**After**

```yaml
apiVersion: neuvector.com/v1
kind: NvSecurityRule
metadata:
  name: nv.bookstore-ui.publishing-company
  namespace: publishing-company
spec:
  egress:
    - action: allow
      applications:
        - Consul
        - SSL
      name: nv.consul-server.consul-egress-0
      ports: any
    - action: allow
      name: nv.consul-server.consul-egress-1
      ports: tcp/8502
  ingress:
    - action: allow
      name: nv.ui-ingress-0
    - action: allow
      name: nv.ui-ingress-1
  process_profile:
    baseline: zero-drift
    mode: Protect
```

- **Root** uses K8s order: `apiVersion`, `kind`, `metadata`, `spec`.
- **metadata** and **spec** keys are alphabetical.
- **spec.egress** and **spec.ingress** are sorted by each item’s `name`.
- Keys inside every mapping remain alphabetical.

---

## Config file reference

| Field  | Meaning                                                                                                                       |
|--------|-------------------------------------------------------------------------------------------------------------------------------|
| `path` | Dot-separated path from the **document root** to the **list** (e.g. `spec.egress`, `spec.ingress`).                           |
| `key`  | For each **element** of that list (a mapping), the field name to sort by (e.g. `name`). Missing keys compare as empty string. |

- You can have as many `listSortKeys` entries as you need (different or nested lists).
- Paths are matched exactly; no wildcards.
- Copy [.yaml-sort.example.yaml](.yaml-sort.example.yaml) to `.yaml-sort.yaml` and adjust paths/keys for your YAML.

See [README](README.md) for installation, flags, and usage.
