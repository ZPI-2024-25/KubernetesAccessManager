Instrukcja instalacji będzie po polsku (może przydać się przy dokumentacji).

# 1. Skrócona instrukcja

```shell
cd charts/
helm install access-manager kam
```
Tyle wystarczy by zainstalować helm charta na klastrze. Więcej wyjaśnień poniżej.

# 2. Komendy helm

```shell
helm install [name] [chart name] --namespace <namespace>
```
Komenda instalująca helm chart na klastrze.
`[name]` odnosi się do tego jak ta konkretna instancja helm chartu będzie nazywana. Wpływa na nazwy zasobów.
`[chart name]` nazwa helm charta (folderu.) W naszym przypadku `kam` z folderu `charts`.
Flaga `--namespace <namespace>` jest opcjonalna i pozwala zainstalować chart w wybranym namespace.

[Więcej informacji](https://helm.sh/docs/helm/helm_install/)

---

```shell
helm uninstall <name>
```
Komenda pozwalająca odinstalować release charta. Po wykonaniu wszystkie zasoby związane z daną instancją charta zostaną usunięte.

[Więcej informacji](https://helm.sh/docs/helm/helm_uninstall/)

---

```shell
helm upgrade <release> <chart>
```
Komenda uaktualniająca wybrany release (instancję charta z określoną nazwą). Wszystkie zmiany wprowadzone w charcie zostaną zastosowane na klastrze. Numer revision zwiększy się o jeden.

[Więcej informacji](https://helm.sh/docs/helm/helm_upgrade/)

---

```shell
helm rollback <release> <revision>
```
Komenda cofająca release do podanej rewizji. Jeżeli revision będzie puste lub równe 0 release zostanie cofnięty do poprzedniej wersji. Rollback tworzy kolejną rewizję (zwiększa revision).

[Więcej informacji](https://helm.sh/docs/helm/helm_rollback/)

---

```shell
helm template [chart name]
```
Komenda renderuje lokalnie szablony i wyświetla jak będą wyglądały. Wartości które normalnie zostałby uzyskane z klastra będą dodane sztucznie. Nie będzie też serwerowych testów (np. czy na klastrze istnieje już taki release).

[Więcej informacji](https://helm.sh/docs/helm/helm_template/)

---

```shell
helm install [name] [chart name] --dry-run --debug
```
Komenda działająca podobnie do `helm template`, z tym że sprawdzi również czy nie ma konfliktujących zasobów na klastrze.
[Więcej infromacji](https://helm.sh/docs/chart_template_guide/debugging/)

---

```shell
helm create <name>
```
Komenda tworząca nowy chart folder wraz z często używanymi plikami.

# 3. Values.yaml

## Ogólne wartości

```yaml
nameOverride: "access-manager"  
fullnameOverride: ""  
  
podAnnotations: {}  
podLabels: {}
```
**nameOverride** - zastępuje nazwę chartu (znajdującą się w Chart.yaml) podczas tworzenia nazw obiektów Kubernetes. Nasza nazwa jest długa, dlatego ją skróciłem.
**fullnameOverride** - całkowicie zastępuje wygenerowaną nazwę

*Uwaga: jeżeli release name oraz chart name (lub name Override) będą takie same, wybrane zostanie release name. Dlatego mając access-manager i access-manager nazwy zawierają jedynie access-manager a nie access-manager-access-manager*

**podAnnotations** - dodaje do podów zaprezentowane adnotacje
**podLabels** - pozwala dodać dodatkowe etykiety do podów. Chart bazowo dodaje własne. Są zdefiniowane w `_helpers.tpl`, dokładniej w `charts.labels` i jego pochodnych

## Backend

```yaml
backend:  
  replicaCount: 1  
  image:  
    repository: ninjashadow/kubernetes-access-manager-backend  
    pullPolicy: IfNotPresent  
    # Overrides the image tag whose default is the chart appVersion.  
    tag: ""  
  
  serviceAccount:  
    create: true  
    # Automatically mount a ServiceAccount's API credentials?  
    automount: true  
    annotations: {}  
    # The name of the service account to use.  
    # If not set and create is true, a name is generated using the fullname template    name: ""  
  
  service:  
    type: ClusterIP  
    port: 8080  
  
  livenessProbe:  
    httpGet:  
      path: /live  
      port: 8082  
  readinessProbe:  
    httpGet:  
      path: /ready  
      port: 8082  
  
  autoscaling:  
    enabled: false  
    minReplicas: 1  
    maxReplicas: 100  
    targetCPUUtilizationPercentage: 80  
    # targetMemoryUtilizationPercentage: 80  
  
  rbac:  
    create: true  
    rules:  
      - apiGroups: [ "*" ]  
        resources: [ "*" ]  
        verbs: [ "get", "list", "create", "update", "delete", "patch" ]
```
**pullPolicy** - określa kiedy Kubernetes próbuje ściągnąć określony obraz
- **IfNotPresent** - tylko jeśli obrazu nie ma lokalnie
- **Always** - szuka obrazu za każdym razem gdy uruchamiany jest kontener
- **Never** - nigdy nie ściąga, korzysta tylko z lokalnego obrazu (jeśli jest)
  **tag** - nadpisuje tag obrazu. Domyślnie tag obrazu równy jest "v" + appVersion z pliku Chart.yaml
  **livenessProbe** - określa na jakiej ścieżce i porcie kubernetes ma szukać live probe
  **readinessProbe** - to samo co livenessProbe
  **autoscaling** - pozwala włączyć opcjonalne skalowanie liczby replik w zależności od obciążenia
  **rbac** - jeżeli klaster używa rbac aplikacja wymaga określonych uprawnień. Ustawienie create na true tworzy `Cluster Role` i `Cluster Role Binding` z określonymi uprawnieniami

## Frontend

To samo co w Backendzie

## Ingress

```yaml
ingress:  
  enabled: true  
  className: "nginx"  
  annotations:  
#     kubernetes.io/ingress.class: nginx  
#     kubernetes.io/tls-acme: "true"  
  
  hosts:  
    - host: "kam.local"  
      paths:  
        - backend:  
            - path: /api/v1/k8s  
              pathType: Prefix  
            - path: /api/v1/auth  
              pathType: Prefix  
            - path: /api/v1/helm  
              pathType: Prefix  
        - frontend:  
            - path: /  
              pathType: Prefix  
  tls: []  
  #  - secretName: chart-example-tls  
  #    hosts:  #      - chart-example.local
```
Opcjonalna konfiguracja sieciowa.
**className** - ingress może korzystać z różnych **Ingress Controllerów**. className służy do wyboru używanego przez ingress controllera. W późniejszej sekcji będzie więcej o minikube i ingress.
**hosts** - lista hostów z własnymi parametrami. Można zdefiniować wiele hostów.
**host** - pozwala mieć wiele hostów na jednym adresie ip. Jedna połączenie do ingress wymaga właściwego hosta. Jeżeli obecny, ścieżki i zasady są stosowane tylko do ich hosta.
**paths** - dostępne dla hosta ścieżki
**backend** - ścieżki dla danego hosta z portem i serwisem backendowym.
**frontend** - ścieżki dla danego hosta z portem i serwisem frontendowym.
**pathType** - typ ścieżki.
- **ImplementationSpecific** - zależy od IngressClass
- **Exact** - dokładne dopasowanie, rozróżnia wielkość liter oraz obecność lub brak /
- **Prefix** - dopasowuje na podsawie ścieżki podzielonej na /. Rozróżnia wielkość liter. Luźniejsze niż Exact
  **tls** - pozwala zabezpieczyć Ingress za pomocą kluczy i certyfikatów. Wymaga Secret na klastrze.
  secretName - nazwa zasobu zawierającego klucze i certyfikaty TLS. Przykład:
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: testsecret-tls
  namespace: default
data:
  tls.crt: base64 encoded cert
  tls.key: base64 encoded key
type: kubernetes.io/tls
```
**hosts** - muszą pasować do hostów w sekcji hosts.

# 4. Połączenie pomiędzy Frontendem i Backendem

Jako metodę przekazywania danych pomiędzy frontendem i backendem wybrałem bezpośrednią łączność w klastrze, bez wychodzenia na zewnątrz. Dzieje się to poprzez konfigurację nginx i wartości środowiskowe.

Fragment `kam/templates/frontend/frontend-dep.yaml`
```yaml
...
env:  
  - name: BACKEND_SERVICE_HOST  
    value: {{ include "charts.fullnameBackend" . }}  
  - name: BACKEND_SERVICE_PORT  
    value: {{ .Values.backend.service.port | quote }}
```

# 5. Testowanie lokalnie

Do postawienia lokalnego klastra Kubernetes wykorzystałem minikube ze sterownikiem docker. Wszystko pokazane poniżej będzie odnosić się do tej konfiguracji.

## Bez Ingress

Domyślnie aplikacja posiada włączony Ingress. Aby go wyłączyć trzeba ustawić `values.yaml` w następujący sposób:
```yaml
ingress:  
  enabled: false
```

Dodatkowo należy zmienić typy serwisów na NodePort albo LoadBalancer.

Więcej informacji czym się różnią [tutaj](https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types)

W minikube można pominąć nawet tą zmianę, gdyż do celów testowych można odsłonić również typ CluserIp.

Aby uzyskać dostęp do serwisu wystarczy wpisać:
```shell
minikube service [service name]
```
Po chwili pojawi się url i zostanie otworzony w przeglądarce.

## Ingress

Minikube posiada addon, który zawiera skonfigurowany Ingress Controller (dokładniej Nginx Controleler).

Można go uruchomić za pomocą:
```shell
minikube addons enable ingress
```

Aby sprawdzić czy wszystko się udało można wpisać:
```shell
kubectl get pods -n ingress-nginx
```

Wynik powinien być podobny do tego:
```
NAME                                        READY   STATUS      RESTARTS    AGE
ingress-nginx-admission-create-g9g49        0/1     Completed   0          11m
ingress-nginx-admission-patch-rqp78         0/1     Completed   1          11m
ingress-nginx-controller-59b45fb494-26npt   1/1     Running     0          11m
```

Jeżeli wszystko działa poprawnie, najlepiej otworzyć inny terminal i wpisać:
```shell
minikube tunnel
```
Otworzy to dostęp do ingress na adresie 127.0.0.1

Z racji tego że korzystamy w ingress z hostów, samo wpisanie http://127.0.0.1 nie przekieruje nas do aplikacji. Należy dodać lokalnego hosta.

### Windows
Aby dodać lokalnego hosta należy edytować plik `C:\Windows\System32\drivers\etc\hosts`.
**Pamiętaj o uprawnieniach administratora!**

### Linux
Aby dodać lokalnego hosta należy edytować plik `/etc/hosts`.