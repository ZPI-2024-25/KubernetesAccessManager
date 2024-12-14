## For English Press [Here]()

# Architektura

## Ogólna struktura projektu
![Diagram komponentów](images/component_diagram.png)

Kubernetes Access Manager (KAM) składa się z trzech głównych komponentów:
- **Frontend**- aplikacja webowa, która pozwala na zarządzanie zasobami w klastrze Kubernetes. 
Frontend jest odpowiedzialny za wyświetlanie interfejsu użytkownika, a także za komunikację z backendem.
Do uwierzytelniania użytkowników wykorzystuje on tokeny JWT.
- **Backend**- serwer aplikacyjny, który zarządza zasobami w klastrze Kubernetes.
Backend jest odpowiedzialny za autoryzację użytkowników, a także za komunikację z API klastra Kubernetes.
Do uwierzytelniania użytkowników wykorzystuje on tokeny JWT.
- **Identity Provider**- dostawca tożsamości, który pozwala na uwierzytelnianie i autoryzację użytkowników.
Wybrany został protokół OpenID Connect, a jako dostawca tożsamości wykorzystany został Keycloak.

Kam komunikuje się z API klastra Kubernetes za pomocą biblioteki `client-go`, natomiast z Klientem Helmowym za pomocą biblioteki `helm`.

## Helm chart
Aplikacja uruchamiana jest w klastrze Kubernetes. W celu łatwej instalacji i zarządzania aplikacją wykorzystany został Helm.

![Diagram Helm Chart](images/helm_chart.png)

Helm chart aplikacji składa się z czterech głównych komponentów:
- **Backend**- serwer aplikacyjny, który zarządza zasobami w klastrze Kubernetes. Składają się na niego:
    - Pod- w nich uruchamiana jest instancja serwera backendowego.
    - Deployment- definiuje, w jaki sposób mają być tworzone i zarządzane instancje podów.
    - Service- umożliwia komunikację pomiędzy częścią frontendową a backendową wewnątrz klastra Kubernetes.
    - Horizontal Pod Autoscaler- umożliwia automatyczne skalowanie liczby podów w zależności od obciążenia.
- **Frontend**- serwer aplikacji webowej, która pozwala na zarządzanie zasobami w klastrze Kubernetes. Składają się na niego:
    - Pod- w nich uruchamiana jest instancja serwera frontendowego.
    - Deployment- definiuje, w jaki sposób mają być tworzone i zarządzane instancje podów.
    - Service- umożliwia komunikację pomiędzy częścią frontendową a backendową wewnątrz klastra Kubernetes.
    - Horizontal Pod Autoscaler- umożliwia automatyczne skalowanie liczby podów w zależności od obciążenia.
- **Access Control**- zbiór zasobów przyznający uprawnienia wymagane w części backendowej. Składają się na niego:
    - Service Account-umożliwia aplikacji dostęp do API klastra Kubernetes.
    - Cluster Role- definiuje zestaw uprawnień wymaganych do zarządzania zasobami w klastrze Kubernetes.
    - Cluster Role Binding- definiuje powiązanie pomiędzy Service Accountem a Cluster Role.
- **Ingress**- definiuje reguły zarządzania ruchem przychodzącym do klastra Kubernetes.

Część **Access Control** jest wymagana tylko w przypadku, gdy na klastrze włączony jest mechanizm RBAC.

## Część backendowa