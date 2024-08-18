# Escenario 1

## Requerimiento

- cluster de kubernetes (preferible con Loadbalancer implementado)
- helm v3 instalado en su maquina local

## Contexto

Una empresa de comercio electrónico mediana quiere expandir su negocio y optimizar operaciones usando Odoo, un ERP de código abierto para gestionar tienda en línea, inventario y CRM. Su equipo de desarrollo, familiarizado con Kubernetes, planea desplegar Odoo en su clúster existente para mejorar la escalabilidad.
Necesitan configurar Odoo con un correo electrónico personalizado y contraseña, cargar datos de prueba para probar, y enviar todos los correos a través de un servicio SMTP interno.

## Solución

1. Busquemos si hay algun artifact disponible en [artifacthub.io](https://artifacthub.io/)
2. Encontraras un helm chart disponible de bitnami/odoo.
3. Revisa los [parametros del Chart](https://artifacthub.io/packages/helm/bitnami/odoo#parameters) y genera un nuevo archivo con los valores necesarios para desplegarlo (en este caso este archivo esta en [app_values.yaml](app_values.yaml))
4. Utiliza los siguientes comandos para desplegar el cluster:

```sh
# Agregar el repositorio de bitnami
$ helm repo add bitnami https://charts.bitnami.com/bitnami

# Asegura que el Chart de Odoo este disponible en los repositorio de bitnami
$ helm search repo odoo

# Despliega el chart en el cluster (Fijate en el output de este comando, tiene informacion relevante sobre la instalación)
$ helm install my-odoo-app -f app_values.yaml bitnami/odoo

# En caso de que estes usando Minikube, puedes abrir otra pestaña para acceder al UI del proyecto
minikube service my-odoo-app --url
```
