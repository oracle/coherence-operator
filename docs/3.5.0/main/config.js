function createConfig() {
    return {
        home: "docs/about/01_overview",
        release: "3.5.0",
        releases: [
            "3.5.0"
        ],
        pathColors: {
            "*": "blue-grey"
        },
        theme: {
            primary: '#1976D2',
            secondary: '#424242',
            accent: '#82B1FF',
            error: '#FF5252',
            info: '#2196F3',
            success: '#4CAF50',
            warning: '#FFC107'
        },
        navTitle: 'Coherence Operator',
        navIcon: null,
        navLogo: 'images/logo.png'
    };
}

function createRoutes(){
    return [
        {
            path: '/docs/about/01_overview',
            meta: {
                h1: 'Overview',
                title: 'Overview',
                h1Prefix: null,
                description: 'Coherence Operator Documentation',
                keywords: 'oracle coherence, kubernetes, operator, documentation',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-about-01_overview', '/docs/about/01_overview', {})
        },
        {
            path: '/docs/about/02_introduction',
            meta: {
                h1: 'Coherence Operator Introduction',
                title: 'Coherence Operator Introduction',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Introduction',
                keywords: 'oracle coherence, kubernetes, operator, documentation, introduction',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-about-02_introduction', '/docs/about/02_introduction', {})
        },
        {
            path: '/docs/about/03_quickstart',
            meta: {
                h1: 'Quick Start',
                title: 'Quick Start',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Quick Start',
                keywords: 'oracle coherence, kubernetes, operator, documentation, quick start',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-about-03_quickstart', '/docs/about/03_quickstart', {})
        },
        {
            path: '/docs/about/04_coherence_spec',
            meta: {
                h1: 'Coherence Operator API Docs',
                title: 'Coherence Operator API Docs',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-about-04_coherence_spec', '/docs/about/04_coherence_spec', {})
        },
        {
            path: '/docs/about/05_upgrade',
            meta: {
                h1: 'Upgrade from Version 2',
                title: 'Upgrade from Version 2',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Upgrade from Version 2',
                keywords: 'oracle coherence, kubernetes, operator, documentation, upgrade, version 2, version 3',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-about-05_upgrade', '/docs/about/05_upgrade', {})
        },
        {
            path: '/docs/installation/001_installation',
            meta: {
                h1: 'Coherence Operator Installation',
                title: 'Coherence Operator Installation',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-installation-001_installation', '/docs/installation/001_installation', {})
        },
        {
            path: '/docs/installation/011_install_manifests',
            meta: {
                h1: 'Install Using Manifests',
                title: 'Install Using Manifests',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-installation-011_install_manifests', '/docs/installation/011_install_manifests', {})
        },
        {
            path: '/docs/installation/012_install_helm',
            meta: {
                h1: 'Install Using Helm',
                title: 'Install Using Helm',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-installation-012_install_helm', '/docs/installation/012_install_helm', {})
        },
        {
            path: '/docs/installation/013_install_kustomize',
            meta: {
                h1: 'Install Using Kustomize',
                title: 'Install Using Kustomize',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-installation-013_install_kustomize', '/docs/installation/013_install_kustomize', {})
        },
        {
            path: '/docs/installation/014_install_openshift',
            meta: {
                h1: 'Install on OpenShift',
                title: 'Install on OpenShift',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-installation-014_install_openshift', '/docs/installation/014_install_openshift', {})
        },
        {
            path: '/docs/installation/015_install_olm',
            meta: {
                h1: 'Install Using OLM',
                title: 'Install Using OLM',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-installation-015_install_olm', '/docs/installation/015_install_olm', {})
        },
        {
            path: '/docs/installation/016_install_tanzu',
            meta: {
                h1: 'Install On Tanzu',
                title: 'Install On Tanzu',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-installation-016_install_tanzu', '/docs/installation/016_install_tanzu', {})
        },
        {
            path: '/docs/installation/020_RBAC',
            meta: {
                h1: 'RBAC Roles',
                title: 'RBAC Roles',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - RBAC Roles',
                keywords: 'oracle coherence, kubernetes, operator, documentation, RBAC, Roles',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-installation-020_RBAC', '/docs/installation/020_RBAC', {})
        },
        {
            path: '/docs/installation/030_pre_release_versions',
            meta: {
                h1: 'Accessing Pre-Release Versions',
                title: 'Accessing Pre-Release Versions',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Accessing Pre-Release Versions',
                keywords: 'oracle coherence, kubernetes, operator, documentation, Accessing Pre-Release Versions',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-installation-030_pre_release_versions', '/docs/installation/030_pre_release_versions', {})
        },
        {
            path: '/docs/installation/040_obtain_coherence_images',
            meta: {
                h1: 'Obtain Coherence Images',
                title: 'Obtain Coherence Images',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Obtain Coherence Images',
                keywords: 'oracle coherence, kubernetes, operator, documentation, Obtain Coherence Images, images',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-installation-040_obtain_coherence_images', '/docs/installation/040_obtain_coherence_images', {})
        },
        {
            path: '/docs/installation/050_private_repos',
            meta: {
                h1: 'Using Private Image Registries',
                title: 'Using Private Image Registries',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Using Private Image Registries',
                keywords: 'oracle coherence, kubernetes, operator, documentation, Private Image Registries, registry',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-installation-050_private_repos', '/docs/installation/050_private_repos', {})
        },
        {
            path: '/docs/installation/070_webhooks',
            meta: {
                h1: 'Operator Web-Hooks',
                title: 'Operator Web-Hooks',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Operator Web-Hooks',
                keywords: 'oracle coherence, kubernetes, operator, documentation, Web-Hooks',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-installation-070_webhooks', '/docs/installation/070_webhooks', {})
        },
        {
            path: '/docs/installation/080_networking',
            meta: {
                h1: 'O/S Networking Configuration',
                title: 'O/S Networking Configuration',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - O/S Networking Configuration',
                keywords: 'oracle coherence, kubernetes, operator, documentation, Networking, Configuration',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-installation-080_networking', '/docs/installation/080_networking', {})
        },
        {
            path: '/docs/applications/010_overview',
            meta: {
                h1: 'Overview',
                title: 'Overview',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Application Deployment',
                keywords: 'oracle coherence, kubernetes, operator, documentation, application deployment, deployment',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-applications-010_overview', '/docs/applications/010_overview', {})
        },
        {
            path: '/docs/applications/020_build_application',
            meta: {
                h1: 'Build Application Images',
                title: 'Build Application Images',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Build Application Images',
                keywords: 'oracle coherence, kubernetes, operator, documentation, Build Application Images, build',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-applications-020_build_application', '/docs/applications/020_build_application', {})
        },
        {
            path: '/docs/applications/030_deploy_application',
            meta: {
                h1: 'Deploy Coherence Applications',
                title: 'Deploy Coherence Applications',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Deploy Coherence Applications',
                keywords: 'oracle coherence, kubernetes, operator, documentation, Deploy Coherence Applications',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-applications-030_deploy_application', '/docs/applications/030_deploy_application', {})
        },
        {
            path: '/docs/applications/032_rolling_upgrade',
            meta: {
                h1: 'Rolling Upgrades of Coherence Applications',
                title: 'Rolling Upgrades of Coherence Applications',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Rolling Upgrades of Coherence Applications',
                keywords: 'oracle coherence, kubernetes, operator, documentation, rolling upgrades',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-applications-032_rolling_upgrade', '/docs/applications/032_rolling_upgrade', {})
        },
        {
            path: '/docs/applications/040_application_main',
            meta: {
                h1: 'Set the Application Main',
                title: 'Set the Application Main',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Set the Application Main',
                keywords: 'oracle coherence, kubernetes, operator, documentation, set application main',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-applications-040_application_main', '/docs/applications/040_application_main', {})
        },
        {
            path: '/docs/applications/050_application_args',
            meta: {
                h1: 'Set Application Arguments',
                title: 'Set Application Arguments',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Set Application Arguments',
                keywords: 'oracle coherence, kubernetes, operator, documentation, set application arguments',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-applications-050_application_args', '/docs/applications/050_application_args', {})
        },
        {
            path: '/docs/applications/060_application_working_dir',
            meta: {
                h1: 'Set the Working Directory',
                title: 'Set the Working Directory',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Set the Working Directory',
                keywords: 'oracle coherence, kubernetes, operator, documentation, Set the Working Directory',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-applications-060_application_working_dir', '/docs/applications/060_application_working_dir', {})
        },
        {
            path: '/docs/applications/070_spring',
            meta: {
                h1: 'Spring Boot Applications',
                title: 'Spring Boot Applications',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Spring Boot Applications',
                keywords: 'oracle coherence, kubernetes, operator, documentation, Spring Boot Applications',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-applications-070_spring', '/docs/applications/070_spring', {})
        },
        {
            path: '/docs/applications/080_entrypoint',
            meta: {
                h1: 'Run an Image Entry Point',
                title: 'Run an Image Entry Point',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Run an Image Entry Point',
                keywords: 'oracle coherence, kubernetes, operator, documentation, Run an Image Entry Point',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-applications-080_entrypoint', '/docs/applications/080_entrypoint', {})
        },
        {
            path: '/docs/coherence/010_overview',
            meta: {
                h1: 'Overview',
                title: 'Overview',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Coherence Settings',
                keywords: 'oracle coherence, kubernetes, operator, documentation, coherence settings',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-coherence-010_overview', '/docs/coherence/010_overview', {})
        },
        {
            path: '/docs/coherence/020_cluster_name',
            meta: {
                h1: 'Coherence Cluster Name',
                title: 'Coherence Cluster Name',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Set Coherence Cluster Name',
                keywords: 'oracle coherence, kubernetes, operator, documentation, set coherence cluster name, cluster name',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-coherence-020_cluster_name', '/docs/coherence/020_cluster_name', {})
        },
        {
            path: '/docs/coherence/021_member_identity',
            meta: {
                h1: 'Member Identity',
                title: 'Member Identity',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Member Identity',
                keywords: 'oracle coherence, kubernetes, operator, documentation, member identity, member',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-coherence-021_member_identity', '/docs/coherence/021_member_identity', {})
        },
        {
            path: '/docs/coherence/030_cache_config',
            meta: {
                h1: 'Cache Configuration File',
                title: 'Cache Configuration File',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Set the Cache Configuration File',
                keywords: 'oracle coherence, kubernetes, operator, documentation, Cache Configuration File, cache configuration',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-coherence-030_cache_config', '/docs/coherence/030_cache_config', {})
        },
        {
            path: '/docs/coherence/040_override_file',
            meta: {
                h1: 'Operational Configuration File',
                title: 'Operational Configuration File',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Set the Operational Configuration File',
                keywords: 'oracle coherence, kubernetes, operator, documentation, Cache Operational Configuration File, operational configuration',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-coherence-040_override_file', '/docs/coherence/040_override_file', {})
        },
        {
            path: '/docs/coherence/050_storage_enabled',
            meta: {
                h1: 'Storage Enabled or Disabled',
                title: 'Storage Enabled or Disabled',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Storage Enabled or Disabled Deployments',
                keywords: 'oracle coherence, kubernetes, operator, documentation, Storage Enabled or Disabled Deployments',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-coherence-050_storage_enabled', '/docs/coherence/050_storage_enabled', {})
        },
        {
            path: '/docs/coherence/060_log_level',
            meta: {
                h1: 'Coherence Log Level',
                title: 'Coherence Log Level',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Set the Coherence Log Level',
                keywords: 'oracle coherence, kubernetes, operator, documentation, Set the Coherence Log Level',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-coherence-060_log_level', '/docs/coherence/060_log_level', {})
        },
        {
            path: '/docs/coherence/070_wka',
            meta: {
                h1: 'Well Known Addressing',
                title: 'Well Known Addressing',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Well Known Addressing and Cluster Discovery',
                keywords: 'oracle coherence, kubernetes, operator, documentation, Well Known Addressing and Cluster Discovery, WKA',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-coherence-070_wka', '/docs/coherence/070_wka', {})
        },
        {
            path: '/docs/coherence/080_persistence',
            meta: {
                h1: 'Coherence Persistence',
                title: 'Coherence Persistence',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Configure Coherence Persistence',
                keywords: 'oracle coherence, kubernetes, operator, documentation, configure Coherence Persistence, persistence',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-coherence-080_persistence', '/docs/coherence/080_persistence', {})
        },
        {
            path: '/docs/coherence/090_ipmonitor',
            meta: {
                h1: 'Coherence IPMonitor',
                title: 'Coherence IPMonitor',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Coherence IPMonitor',
                keywords: 'oracle coherence, kubernetes, operator, documentation, IPMonitor',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-coherence-090_ipmonitor', '/docs/coherence/090_ipmonitor', {})
        },
        {
            path: '/docs/jvm/010_overview',
            meta: {
                h1: 'Overview',
                title: 'Overview',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - JVM Settings',
                keywords: 'oracle coherence, kubernetes, operator, documentation, JVM, Settings',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-jvm-010_overview', '/docs/jvm/010_overview', {})
        },
        {
            path: '/docs/jvm/020_classpath',
            meta: {
                h1: 'Set the Classpath',
                title: 'Set the Classpath',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Set the Classpath',
                keywords: 'oracle coherence, kubernetes, operator, documentation, Classpath',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-jvm-020_classpath', '/docs/jvm/020_classpath', {})
        },
        {
            path: '/docs/jvm/030_jvm_args',
            meta: {
                h1: 'Arbitrary JVM Arguments',
                title: 'Arbitrary JVM Arguments',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Arbitrary JVM Arguments',
                keywords: 'oracle coherence, kubernetes, operator, documentation, Arbitrary, JVM, Arguments',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-jvm-030_jvm_args', '/docs/jvm/030_jvm_args', {})
        },
        {
            path: '/docs/jvm/040_gc',
            meta: {
                h1: 'Garbage Collector Settings',
                title: 'Garbage Collector Settings',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Garbage Collector Settings',
                keywords: 'oracle coherence, kubernetes, operator, documentation, Garbage Collector',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-jvm-040_gc', '/docs/jvm/040_gc', {})
        },
        {
            path: '/docs/jvm/050_memory',
            meta: {
                h1: 'Heap & Memory Settings',
                title: 'Heap & Memory Settings',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Heap &amp; Memory Settings',
                keywords: 'oracle coherence, kubernetes, operator, documentation, Heap, Memory, Settings',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-jvm-050_memory', '/docs/jvm/050_memory', {})
        },
        {
            path: '/docs/jvm/070_debugger',
            meta: {
                h1: 'Debugger Configuration',
                title: 'Debugger Configuration',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Debugger Configuration',
                keywords: 'oracle coherence, kubernetes, operator, documentation, Debugger, Configuration',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-jvm-070_debugger', '/docs/jvm/070_debugger', {})
        },
        {
            path: '/docs/jvm/090_container_limits',
            meta: {
                h1: 'Container Resource Limits',
                title: 'Container Resource Limits',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Container Resource Limits',
                keywords: 'oracle coherence, kubernetes, operator, documentation, Container, Resource Limits',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-jvm-090_container_limits', '/docs/jvm/090_container_limits', {})
        },
        {
            path: '/docs/ports/010_overview',
            meta: {
                h1: 'Overview',
                title: 'Overview',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Ports Overview',
                keywords: 'oracle coherence, kubernetes, operator, ports',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-ports-010_overview', '/docs/ports/010_overview', {})
        },
        {
            path: '/docs/ports/020_container_ports',
            meta: {
                h1: 'Additional Container Ports',
                title: 'Additional Container Ports',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Additional Container Ports',
                keywords: 'oracle coherence, kubernetes, operator, additional container ports, container ports, ports',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-ports-020_container_ports', '/docs/ports/020_container_ports', {})
        },
        {
            path: '/docs/ports/030_services',
            meta: {
                h1: 'Configure Services for Ports',
                title: 'Configure Services for Ports',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Configure Services for Ports',
                keywords: 'oracle coherence, kubernetes, operator, services, ports',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-ports-030_services', '/docs/ports/030_services', {})
        },
        {
            path: '/docs/ports/040_servicemonitors',
            meta: {
                h1: 'Prometheus ServiceMonitors',
                title: 'Prometheus ServiceMonitors',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Prometheus ServiceMonitors',
                keywords: 'oracle coherence, kubernetes, operator, prometheus, serviceMonitors',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-ports-040_servicemonitors', '/docs/ports/040_servicemonitors', {})
        },
        {
            path: '/docs/networking/010_overview',
            meta: {
                h1: 'Overview',
                title: 'Overview',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Networking Overview',
                keywords: 'oracle coherence, kubernetes, operator, networking, network',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-networking-010_overview', '/docs/networking/010_overview', {})
        },
        {
            path: '/docs/networking/020_dual_stack',
            meta: {
                h1: 'Dual Stack Networking',
                title: 'Dual Stack Networking',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Dual Stack Networking',
                keywords: 'oracle coherence, kubernetes, operator, networking, dual stack',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-networking-020_dual_stack', '/docs/networking/020_dual_stack', {})
        },
        {
            path: '/docs/scaling/010_overview',
            meta: {
                h1: 'Scale Coherence Deployments',
                title: 'Scale Coherence Deployments',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Scale Coherence Deployments',
                keywords: 'oracle coherence, kubernetes, operator, scale coehrence, scale deployments',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-scaling-010_overview', '/docs/scaling/010_overview', {})
        },
        {
            path: '/docs/ordering/010_overview',
            meta: {
                h1: 'Deployment Start Order',
                title: 'Deployment Start Order',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Deployment Start Order',
                keywords: 'oracle coherence, kubernetes, operator, deployment, start order, ordering',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-ordering-010_overview', '/docs/ordering/010_overview', {})
        },
        {
            path: '/docs/management/010_overview',
            meta: {
                h1: 'Overview',
                title: 'Overview',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Management Overview',
                keywords: 'oracle coherence, kubernetes, operator, documentation, management',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-management-010_overview', '/docs/management/010_overview', {})
        },
        {
            path: '/docs/management/020_management_over_rest',
            meta: {
                h1: 'Management over REST',
                title: 'Management over REST',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Management over REST',
                keywords: 'oracle coherence, kubernetes, operator, Management, REST',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-management-020_management_over_rest', '/docs/management/020_management_over_rest', {})
        },
        {
            path: '/docs/management/025_coherence_cli',
            meta: {
                h1: 'The Coherence CLI',
                title: 'The Coherence CLI',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Coherence CLI',
                keywords: 'oracle coherence, kubernetes, operator, coherence-cli, cli',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-management-025_coherence_cli', '/docs/management/025_coherence_cli', {})
        },
        {
            path: '/docs/management/026_queryplus',
            meta: {
                h1: 'The Coherence Query Plus',
                title: 'The Coherence Query Plus',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Coherence Query Plus',
                keywords: 'oracle coherence, kubernetes, operator, Query Plus, cohql',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-management-026_queryplus', '/docs/management/026_queryplus', {})
        },
        {
            path: '/docs/management/027_jshell',
            meta: {
                h1: 'Using JShell',
                title: 'Using JShell',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Using JShell',
                keywords: 'oracle coherence, kubernetes, operator, management, jshell',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-management-027_jshell', '/docs/management/027_jshell', {})
        },
        {
            path: '/docs/management/040_ssl',
            meta: {
                h1: 'SSL with Management over REST',
                title: 'SSL with Management over REST',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - SSL with Management over REST',
                keywords: 'oracle coherence, kubernetes, operator, management, REST, SSL, TLS',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-management-040_ssl', '/docs/management/040_ssl', {})
        },
        {
            path: '/docs/management/100_tmb_test',
            meta: {
                h1: 'Coherence Network Testing',
                title: 'Coherence Network Testing',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Coherence Network Testing',
                keywords: 'oracle coherence, kubernetes, operator, management, network test',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-management-100_tmb_test', '/docs/management/100_tmb_test', {})
        },
        {
            path: '/docs/metrics/010_overview',
            meta: {
                h1: 'Overview',
                title: 'Overview',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Metrics Overview',
                keywords: 'oracle coherence, kubernetes, operator, metrics, overview',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-metrics-010_overview', '/docs/metrics/010_overview', {})
        },
        {
            path: '/docs/metrics/020_metrics',
            meta: {
                h1: 'Publish Metrics',
                title: 'Publish Metrics',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Publish Metrics',
                keywords: 'oracle coherence, kubernetes, operator, metrics, publish',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-metrics-020_metrics', '/docs/metrics/020_metrics', {})
        },
        {
            path: '/docs/metrics/030_importing',
            meta: {
                h1: 'Import the Grafana Dashboards',
                title: 'Import the Grafana Dashboards',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Import the Grafana Dashboards',
                keywords: 'oracle coherence, kubernetes, operator, metrics, grafana, dashboards, import',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-metrics-030_importing', '/docs/metrics/030_importing', {})
        },
        {
            path: '/docs/metrics/040_dashboards',
            meta: {
                h1: 'Coherence Grafana Dashboards',
                title: 'Coherence Grafana Dashboards',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Coherence Grafana Dashboards',
                keywords: 'oracle coherence, kubernetes, operator, metrics, grafana, dashboards',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-metrics-040_dashboards', '/docs/metrics/040_dashboards', {})
        },
        {
            path: '/docs/metrics/050_ssl',
            meta: {
                h1: 'SSL with Metrics',
                title: 'SSL with Metrics',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - SSL with Metrics',
                keywords: 'oracle coherence, kubernetes, operator, metrics, ssl, tls',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-metrics-050_ssl', '/docs/metrics/050_ssl', {})
        },
        {
            path: '/docs/logging/010_overview',
            meta: {
                h1: 'Overview',
                title: 'Overview',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Logging Overview',
                keywords: 'oracle coherence, kubernetes, operator, documentation, Logging, Overview',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-logging-010_overview', '/docs/logging/010_overview', {})
        },
        {
            path: '/docs/logging/020_logging',
            meta: {
                h1: 'Log Capture with Fluentd',
                title: 'Log Capture with Fluentd',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Log Capture with Fluentd',
                keywords: 'oracle coherence, kubernetes, operator, documentation, Log Capture, Fluentd',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-logging-020_logging', '/docs/logging/020_logging', {})
        },
        {
            path: '/docs/logging/030_kibana',
            meta: {
                h1: 'Using Kibana Dashboards',
                title: 'Using Kibana Dashboards',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Using Kibana Dashboards',
                keywords: 'oracle coherence, kubernetes, operator, documentation, Kibana, Dashboards',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-logging-030_kibana', '/docs/logging/030_kibana', {})
        },
        {
            path: '/docs/other/010_overview',
            meta: {
                h1: 'Overview',
                title: 'Overview',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Other Settings',
                keywords: 'oracle coherence, kubernetes, operator, settings',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-other-010_overview', '/docs/other/010_overview', {})
        },
        {
            path: '/docs/other/020_environment',
            meta: {
                h1: 'Environment Variables',
                title: 'Environment Variables',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Environment Variables',
                keywords: 'oracle coherence, kubernetes, operator, Environment Variables',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-other-020_environment', '/docs/other/020_environment', {})
        },
        {
            path: '/docs/other/030_labels',
            meta: {
                h1: 'Pod Labels',
                title: 'Pod Labels',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Pod Labels',
                keywords: 'oracle coherence, kubernetes, operator, pod, labels',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-other-030_labels', '/docs/other/030_labels', {})
        },
        {
            path: '/docs/other/040_annotations',
            meta: {
                h1: 'Adding Annotations',
                title: 'Adding Annotations',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Adding Annotations',
                keywords: 'oracle coherence, kubernetes, operator, annotations',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-other-040_annotations', '/docs/other/040_annotations', {})
        },
        {
            path: '/docs/other/041_global_labels',
            meta: {
                h1: 'Global Labels and Annotations',
                title: 'Global Labels and Annotations',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Global Labels and Annotations',
                keywords: 'oracle coherence, kubernetes, operator, annotations, global labels',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-other-041_global_labels', '/docs/other/041_global_labels', {})
        },
        {
            path: '/docs/other/045_security_context',
            meta: {
                h1: 'Pod & Container SecurityContext',
                title: 'Pod & Container SecurityContext',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Pod &amp; Container SecurityContexts',
                keywords: 'oracle coherence, kubernetes, operator, pod sercurtyContext, container sercurtyContext',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-other-045_security_context', '/docs/other/045_security_context', {})
        },
        {
            path: '/docs/other/050_configmap_volumes',
            meta: {
                h1: 'Add ConfigMap Volumes',
                title: 'Add ConfigMap Volumes',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Add ConfigMap Volumes',
                keywords: 'oracle coherence, kubernetes, operator, configmap, volumes',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-other-050_configmap_volumes', '/docs/other/050_configmap_volumes', {})
        },
        {
            path: '/docs/other/060_secret_volumes',
            meta: {
                h1: 'Add Secrets Volumes',
                title: 'Add Secrets Volumes',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Add Secrets Volumes',
                keywords: 'oracle coherence, kubernetes, operator, secrets, volumes',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-other-060_secret_volumes', '/docs/other/060_secret_volumes', {})
        },
        {
            path: '/docs/other/070_add_volumes',
            meta: {
                h1: 'Add Pod Volumes',
                title: 'Add Pod Volumes',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Add Pod Volumes',
                keywords: 'oracle coherence, kubernetes, operator, pod, volumes',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-other-070_add_volumes', '/docs/other/070_add_volumes', {})
        },
        {
            path: '/docs/other/080_add_containers',
            meta: {
                h1: 'Configure Additional Containers',
                title: 'Configure Additional Containers',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Configure Additional Containers',
                keywords: 'oracle coherence, kubernetes, operator, additional containers',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-other-080_add_containers', '/docs/other/080_add_containers', {})
        },
        {
            path: '/docs/other/090_pod_scheduling',
            meta: {
                h1: 'Configure Pod Scheduling',
                title: 'Configure Pod Scheduling',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Configure Pod Scheduling',
                keywords: 'oracle coherence, kubernetes, operator, pod scheduling, scheduling',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-other-090_pod_scheduling', '/docs/other/090_pod_scheduling', {})
        },
        {
            path: '/docs/other/100_resources',
            meta: {
                h1: 'Container Resource Limits',
                title: 'Container Resource Limits',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Container Resource Limits',
                keywords: 'oracle coherence, kubernetes, operator, container, resource limits',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-other-100_resources', '/docs/other/100_resources', {})
        },
        {
            path: '/docs/other/110_readiness',
            meta: {
                h1: 'Readiness & Liveness Probes',
                title: 'Readiness & Liveness Probes',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Readiness &amp; Liveness Probes',
                keywords: 'oracle coherence, kubernetes, operator, Readiness, Liveness, probes',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-other-110_readiness', '/docs/other/110_readiness', {})
        },
        {
            path: '/examples/000_overview',
            meta: {
                h1: 'Overview',
                title: 'Overview',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('examples-000_overview', '/examples/000_overview', {})
        },
        {
            path: '/examples/015_simple_image/README',
            meta: {
                h1: 'Example Coherence Image using JIB',
                title: 'Example Coherence Image using JIB',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Example Coherence Image using JIB',
                keywords: 'oracle coherence, kubernetes, operator, image, jib',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('examples-015_simple_image-README', '/examples/015_simple_image/README', {})
        },
        {
            path: '/examples/016_simple_docker_image/README',
            meta: {
                h1: 'Example Coherence Image using a Dockerfile',
                title: 'Example Coherence Image using a Dockerfile',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Example Coherence Image using a Dockerfile',
                keywords: 'oracle coherence, kubernetes, operator, image, docker, dockerfile',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('examples-016_simple_docker_image-README', '/examples/016_simple_docker_image/README', {})
        },
        {
            path: '/examples/020_hello_world/README',
            meta: {
                h1: 'A \"Hello World\" Operator Example',
                title: 'A \"Hello World\" Operator Example',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - A \"Hello World\" Operator Example',
                keywords: 'oracle coherence, kubernetes, operator, hello world, example',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('examples-020_hello_world-README', '/examples/020_hello_world/README', {})
        },
        {
            path: '/examples/021_deployment/README',
            meta: {
                h1: 'Coherence Deployment Example',
                title: 'Coherence Deployment Example',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Coherence Deployment Example',
                keywords: 'oracle coherence, kubernetes, operator, deployment, example',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('examples-021_deployment-README', '/examples/021_deployment/README', {})
        },
        {
            path: '/examples/025_extend_client/README',
            meta: {
                h1: 'Coherence Extend Clients',
                title: 'Coherence Extend Clients',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Coherence Extend Clients Example',
                keywords: 'oracle coherence, kubernetes, operator, extend client, coherence extend, example',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('examples-025_extend_client-README', '/examples/025_extend_client/README', {})
        },
        {
            path: '/examples/090_tls/README',
            meta: {
                h1: 'Secure Coherence Using TLS',
                title: 'Secure Coherence Using TLS',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Secure Coherence Using TLS Example',
                keywords: 'oracle coherence, kubernetes, operator, secure, tls, ssl, example',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('examples-090_tls-README', '/examples/090_tls/README', {})
        },
        {
            path: '/examples/095_network_policies/README',
            meta: {
                h1: 'Using Network Policies',
                title: 'Using Network Policies',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Using Network Policies Example',
                keywords: 'oracle coherence, kubernetes, operator, network, policies, example',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('examples-095_network_policies-README', '/examples/095_network_policies/README', {})
        },
        {
            path: '/examples/100_federation/README',
            meta: {
                h1: 'Coherence Federation',
                title: 'Coherence Federation',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Coherence Federation Example',
                keywords: 'oracle coherence, kubernetes, operator, federation, OCI, Oracle Cloud Infrastructure, example',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('examples-100_federation-README', '/examples/100_federation/README', {})
        },
        {
            path: '/examples/200_autoscaler/README',
            meta: {
                h1: 'Autoscaling Coherence Clusters',
                title: 'Autoscaling Coherence Clusters',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Autoscaling Coherence Clusters Example',
                keywords: 'oracle coherence, kubernetes, operator, autoscaler, example',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('examples-200_autoscaler-README', '/examples/200_autoscaler/README', {})
        },
        {
            path: '/examples/300_helm/README',
            meta: {
                h1: 'Manage Coherence using Helm',
                title: 'Manage Coherence using Helm',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Manage Coherence using Helm Example',
                keywords: 'oracle coherence, kubernetes, operator, helm, manage, example',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('examples-300_helm-README', '/examples/300_helm/README', {})
        },
        {
            path: '/examples/400_Istio/README',
            meta: {
                h1: 'Using Coherence with Istio',
                title: 'Using Coherence with Istio',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Coherence with Istio Example',
                keywords: 'oracle coherence, kubernetes, operator, istio, example',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('examples-400_Istio-README', '/examples/400_Istio/README', {})
        },
        {
            path: '/examples/900_demo/README',
            meta: {
                h1: 'The Coherence Demo App',
                title: 'The Coherence Demo App',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Coherence Demo Application Example',
                keywords: 'oracle coherence, kubernetes, operator, demo, coherence-demo, example',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('examples-900_demo-README', '/examples/900_demo/README', {})
        },
        {
            path: '/examples/910_polyglot_demo/README',
            meta: {
                h1: 'Polyglot Client Demo',
                title: 'Polyglot Client Demo',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Polyglot Client Demo Example',
                keywords: 'oracle coherence, kubernetes, operator, polyglot, grpc, python, javascript, golang, go, example',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('examples-910_polyglot_demo-README', '/examples/910_polyglot_demo/README', {})
        },
        {
            path: '/examples/no-operator/000_overview',
            meta: {
                h1: 'Overview',
                title: 'Overview',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('examples-no-operator-000_overview', '/examples/no-operator/000_overview', {})
        },
        {
            path: '/examples/no-operator/01_simple_server/README',
            meta: {
                h1: 'A Simple Coherence Cluster',
                title: 'A Simple Coherence Cluster',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - A Simple Coherence Cluster Without the Operator',
                keywords: 'oracle coherence, kubernetes, operator, simple cluster, without operator',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('examples-no-operator-01_simple_server-README', '/examples/no-operator/01_simple_server/README', {})
        },
        {
            path: '/examples/no-operator/02_metrics/README',
            meta: {
                h1: 'Enabling Coherence Metrics',
                title: 'Enabling Coherence Metrics',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Enabling Coherence Metrics Without the Operator',
                keywords: 'oracle coherence, kubernetes, operator, metrics, without operator',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('examples-no-operator-02_metrics-README', '/examples/no-operator/02_metrics/README', {})
        },
        {
            path: '/examples/no-operator/03_extend_tls/README',
            meta: {
                h1: 'Secure Coherence Extend with TLS',
                title: 'Secure Coherence Extend with TLS',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Secure Coherence Extend with TLS Without the Operator',
                keywords: 'oracle coherence, kubernetes, operator, without operator, secure, tls, SSL',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('examples-no-operator-03_extend_tls-README', '/examples/no-operator/03_extend_tls/README', {})
        },
        {
            path: '/examples/no-operator/04_istio/README',
            meta: {
                h1: 'Running Coherence with Istio',
                title: 'Running Coherence with Istio',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Running Coherence with Istio Without the Operator',
                keywords: 'oracle coherence, kubernetes, operator, without operator, istio',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('examples-no-operator-04_istio-README', '/examples/no-operator/04_istio/README', {})
        },
        {
            path: '/examples/no-operator/test-client/README',
            meta: {
                h1: 'Example Extend Client',
                title: 'Example Extend Client',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Example Extend ClientWithout the Operator',
                keywords: 'oracle coherence, kubernetes, operator, without operator, extend client',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('examples-no-operator-test-client-README', '/examples/no-operator/test-client/README', {})
        },
        {
            path: '/docs/troubleshooting/01_trouble-shooting',
            meta: {
                h1: 'Troubleshooting Guide',
                title: 'Troubleshooting Guide',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Troubleshooting Guide',
                keywords: 'oracle coherence, kubernetes, operator, Troubleshooting Guide',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-troubleshooting-01_trouble-shooting', '/docs/troubleshooting/01_trouble-shooting', {})
        },
        {
            path: '/docs/troubleshooting/02_heap_dump',
            meta: {
                h1: 'Capture Heap Dumps',
                title: 'Capture Heap Dumps',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Capture Heap Dumps',
                keywords: 'oracle coherence, kubernetes, operator, catpure, heap dumps',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('docs-troubleshooting-02_heap_dump', '/docs/troubleshooting/02_heap_dump', {})
        },
        {
            path: '/docs/performance/010_performance',
            meta: {
                h1: 'Performance Testing',
                title: 'Performance Testing',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Performance Testing',
                keywords: 'oracle coherence, kubernetes, operator, Performance Testing',
                customLayout: null,
                hasNav: false
            },
            component: loadPage('docs-performance-010_performance', '/docs/performance/010_performance', {})
        },
        {
            path: '/docs/webhooks/01_introduction',
            meta: {
                h1: 'Operator K8s Webhooks',
                title: 'Operator K8s Webhooks',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Operator K8s Webhooks',
                keywords: 'oracle coherence, kubernetes, operator, K8s Webhooks, webhooks',
                customLayout: null,
                hasNav: false
            },
            component: loadPage('docs-webhooks-01_introduction', '/docs/webhooks/01_introduction', {})
        },
        {
            path: '/examples/README',
            meta: {
                h1: 'Examples Overview',
                title: 'Examples Overview',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: false
            },
            component: loadPage('examples-README', '/examples/README', {})
        },
        {
            path: '/examples/no-operator/README',
            meta: {
                h1: 'Coherence in Kubernetes Without the Operator',
                title: 'Coherence in Kubernetes Without the Operator',
                h1Prefix: null,
                description: 'Coherence Operator Documentation - Coherence in Kubernetes Without the Operator',
                keywords: 'oracle coherence, kubernetes, operator, without operator, example',
                customLayout: null,
                hasNav: false
            },
            component: loadPage('examples-no-operator-README', '/examples/no-operator/README', {})
        },
        {
            path: '/', redirect: '/docs/about/01_overview'
        },
        {
            path: '*', redirect: '/'
        }
    ];
}

function createNav(){
    return [
        { header: 'Core documentation' },
        {
            title: 'About',
            action: 'assistant',
            group: '/about',
            items: [
                { href: '/docs/about/01_overview', title: 'Overview' },
                { href: '/docs/about/02_introduction', title: 'Coherence Operator Introduction' },
                { href: '/docs/about/03_quickstart', title: 'Quick Start' },
                { href: '/docs/about/04_coherence_spec', title: 'Coherence Operator API Docs' },
                { href: '/docs/about/05_upgrade', title: 'Upgrade from Version 2' }
            ]
        },
        {
            title: 'Installation',
            action: 'fa-save',
            group: '/install',
            items: [
                { href: '/docs/installation/001_installation', title: 'Coherence Operator Installation' },
                { href: '/docs/installation/011_install_manifests', title: 'Install Using Manifests' },
                { href: '/docs/installation/012_install_helm', title: 'Install Using Helm' },
                { href: '/docs/installation/013_install_kustomize', title: 'Install Using Kustomize' },
                { href: '/docs/installation/014_install_openshift', title: 'Install on OpenShift' },
                { href: '/docs/installation/015_install_olm', title: 'Install Using OLM' },
                { href: '/docs/installation/016_install_tanzu', title: 'Install On Tanzu' },
                { href: '/docs/installation/020_RBAC', title: 'RBAC Roles' },
                { href: '/docs/installation/030_pre_release_versions', title: 'Accessing Pre-Release Versions' },
                { href: '/docs/installation/040_obtain_coherence_images', title: 'Obtain Coherence Images' },
                { href: '/docs/installation/050_private_repos', title: 'Using Private Image Registries' },
                { href: '/docs/installation/070_webhooks', title: 'Operator Web-Hooks' },
                { href: '/docs/installation/080_networking', title: 'O/S Networking Configuration' }
            ]
        },
        {
            title: 'Deploy Applications',
            action: 'cloud_upload',
            group: '/applications',
            items: [
                { href: '/docs/applications/010_overview', title: 'Overview' },
                { href: '/docs/applications/020_build_application', title: 'Build Application Images' },
                { href: '/docs/applications/030_deploy_application', title: 'Deploy Coherence Applications' },
                { href: '/docs/applications/032_rolling_upgrade', title: 'Rolling Upgrades of Coherence Applications' },
                { href: '/docs/applications/040_application_main', title: 'Set the Application Main' },
                { href: '/docs/applications/050_application_args', title: 'Set Application Arguments' },
                { href: '/docs/applications/060_application_working_dir', title: 'Set the Working Directory' },
                { href: '/docs/applications/070_spring', title: 'Spring Boot Applications' },
                { href: '/docs/applications/080_entrypoint', title: 'Run an Image Entry Point' }
            ]
        },
        {
            title: 'Coherence Settings',
            action: 'fa-cogs',
            group: '/coherence',
            items: [
                { href: '/docs/coherence/010_overview', title: 'Overview' },
                { href: '/docs/coherence/020_cluster_name', title: 'Coherence Cluster Name' },
                { href: '/docs/coherence/021_member_identity', title: 'Member Identity' },
                { href: '/docs/coherence/030_cache_config', title: 'Cache Configuration File' },
                { href: '/docs/coherence/040_override_file', title: 'Operational Configuration File' },
                { href: '/docs/coherence/050_storage_enabled', title: 'Storage Enabled or Disabled' },
                { href: '/docs/coherence/060_log_level', title: 'Coherence Log Level' },
                { href: '/docs/coherence/070_wka', title: 'Well Known Addressing' },
                { href: '/docs/coherence/080_persistence', title: 'Coherence Persistence' },
                { href: '/docs/coherence/090_ipmonitor', title: 'Coherence IPMonitor' }
            ]
        },
        {
            title: 'JVM Settings',
            action: 'fa-cog',
            group: '/jvm',
            items: [
                { href: '/docs/jvm/010_overview', title: 'Overview' },
                { href: '/docs/jvm/020_classpath', title: 'Set the Classpath' },
                { href: '/docs/jvm/030_jvm_args', title: 'Arbitrary JVM Arguments' },
                { href: '/docs/jvm/040_gc', title: 'Garbage Collector Settings' },
                { href: '/docs/jvm/050_memory', title: 'Heap & Memory Settings' },
                { href: '/docs/jvm/070_debugger', title: 'Debugger Configuration' },
                { href: '/docs/jvm/090_container_limits', title: 'Container Resource Limits' }
            ]
        },
        {
            title: 'Expose Ports & Services',
            action: 'control_camera',
            group: '/ports',
            items: [
                { href: '/docs/ports/010_overview', title: 'Overview' },
                { href: '/docs/ports/020_container_ports', title: 'Additional Container Ports' },
                { href: '/docs/ports/030_services', title: 'Configure Services for Ports' },
                { href: '/docs/ports/040_servicemonitors', title: 'Prometheus ServiceMonitors' }
            ]
        },
        {
            title: 'Networking',
            action: 'share',
            group: '/networking',
            items: [
                { href: '/docs/networking/010_overview', title: 'Overview' },
                { href: '/docs/networking/020_dual_stack', title: 'Dual Stack Networking' }
            ]
        },
        {
            title: 'Scaling Up & Down',
            action: 'fa-balance-scale',
            group: '/scaling',
            items: [
                { href: '/docs/scaling/010_overview', title: 'Scale Coherence Deployments' }
            ]
        },
        {
            title: 'Start-up Order',
            action: 'line_weight',
            group: '/ordering',
            items: [
                { href: '/docs/ordering/010_overview', title: 'Deployment Start Order' }
            ]
        },
        {
            title: 'Management Diagnostics',
            action: 'fa-stethoscope',
            group: '/management',
            items: [
                { href: '/docs/management/010_overview', title: 'Overview' },
                { href: '/docs/management/020_management_over_rest', title: 'Management over REST' },
                { href: '/docs/management/025_coherence_cli', title: 'The Coherence CLI' },
                { href: '/docs/management/026_queryplus', title: 'The Coherence Query Plus' },
                { href: '/docs/management/027_jshell', title: 'Using JShell' },
                { href: '/docs/management/040_ssl', title: 'SSL with Management over REST' },
                { href: '/docs/management/100_tmb_test', title: 'Coherence Network Testing' }
            ]
        },
        {
            title: 'Metrics',
            action: 'speed',
            group: '/metrics',
            items: [
                { href: '/docs/metrics/010_overview', title: 'Overview' },
                { href: '/docs/metrics/020_metrics', title: 'Publish Metrics' },
                { href: '/docs/metrics/030_importing', title: 'Import the Grafana Dashboards' },
                { href: '/docs/metrics/040_dashboards', title: 'Coherence Grafana Dashboards' },
                { href: '/docs/metrics/050_ssl', title: 'SSL with Metrics' }
            ]
        },
        {
            title: 'Logging',
            action: 'find_in_page',
            group: '/logging',
            items: [
                { href: '/docs/logging/010_overview', title: 'Overview' },
                { href: '/docs/logging/020_logging', title: 'Log Capture with Fluentd' },
                { href: '/docs/logging/030_kibana', title: 'Using Kibana Dashboards' }
            ]
        },
        {
            title: 'Other Settings',
            action: 'widgets',
            group: '/other',
            items: [
                { href: '/docs/other/010_overview', title: 'Overview' },
                { href: '/docs/other/020_environment', title: 'Environment Variables' },
                { href: '/docs/other/030_labels', title: 'Pod Labels' },
                { href: '/docs/other/040_annotations', title: 'Adding Annotations' },
                { href: '/docs/other/041_global_labels', title: 'Global Labels and Annotations' },
                { href: '/docs/other/045_security_context', title: 'Pod & Container SecurityContext' },
                { href: '/docs/other/050_configmap_volumes', title: 'Add ConfigMap Volumes' },
                { href: '/docs/other/060_secret_volumes', title: 'Add Secrets Volumes' },
                { href: '/docs/other/070_add_volumes', title: 'Add Pod Volumes' },
                { href: '/docs/other/080_add_containers', title: 'Configure Additional Containers' },
                { href: '/docs/other/090_pod_scheduling', title: 'Configure Pod Scheduling' },
                { href: '/docs/other/100_resources', title: 'Container Resource Limits' },
                { href: '/docs/other/110_readiness', title: 'Readiness & Liveness Probes' }
            ]
        },
        {
            title: 'Examples',
            action: 'explore',
            group: '/examples',
            items: [
                { href: '/examples/000_overview', title: 'Overview' },
                { href: '/examples/015_simple_image/README', title: 'Example Coherence Image using JIB' },
                { href: '/examples/016_simple_docker_image/README', title: 'Example Coherence Image using a Dockerfile' },
                { href: '/examples/020_hello_world/README', title: 'A \"Hello World\" Operator Example' },
                { href: '/examples/021_deployment/README', title: 'Coherence Deployment Example' },
                { href: '/examples/025_extend_client/README', title: 'Coherence Extend Clients' },
                { href: '/examples/090_tls/README', title: 'Secure Coherence Using TLS' },
                { href: '/examples/095_network_policies/README', title: 'Using Network Policies' },
                { href: '/examples/100_federation/README', title: 'Coherence Federation' },
                { href: '/examples/200_autoscaler/README', title: 'Autoscaling Coherence Clusters' },
                { href: '/examples/300_helm/README', title: 'Manage Coherence using Helm' },
                { href: '/examples/400_Istio/README', title: 'Using Coherence with Istio' },
                { href: '/examples/900_demo/README', title: 'The Coherence Demo App' },
                { href: '/examples/910_polyglot_demo/README', title: 'Polyglot Client Demo' }
            ]
        },
        {
            title: 'Non-Operator Examples',
            action: 'fa-ban',
            group: '/no-operator',
            items: [
                { href: '/examples/no-operator/000_overview', title: 'Overview' },
                { href: '/examples/no-operator/01_simple_server/README', title: 'A Simple Coherence Cluster' },
                { href: '/examples/no-operator/02_metrics/README', title: 'Enabling Coherence Metrics' },
                { href: '/examples/no-operator/03_extend_tls/README', title: 'Secure Coherence Extend with TLS' },
                { href: '/examples/no-operator/04_istio/README', title: 'Running Coherence with Istio' },
                { href: '/examples/no-operator/test-client/README', title: 'Example Extend Client' }
            ]
        },
        {
            title: 'Troubleshooting',
            action: 'fa-question-circle',
            group: '/troubleshooting',
            items: [
                { href: '/docs/troubleshooting/01_trouble-shooting', title: 'Troubleshooting Guide' },
                { href: '/docs/troubleshooting/02_heap_dump', title: 'Capture Heap Dumps' }
            ]
        },
        { divider: true },
        { header: 'Additional resources' },
        {
            title: 'Slack',
            action: 'fa-slack',
            href: 'https://join.slack.com/t/oraclecoherence/shared_invite/enQtNzcxNTQwMTAzNjE4LTJkZWI5ZDkzNGEzOTllZDgwZDU3NGM2YjY5YWYwMzM3ODdkNTU2NmNmNDFhOWIxMDZlNjg2MzE3NmMxZWMxMWE',
            target: '_blank'
        },
        {
            title: 'Coherence Community',
            action: 'people',
            href: 'https://coherence.community',
            target: '_blank'
        },
        {
            title: 'GitHub',
            action: 'fa-github-square',
            href: 'https://github.com/oracle/coherence-operator',
            target: '_blank'
        }
    ];
}