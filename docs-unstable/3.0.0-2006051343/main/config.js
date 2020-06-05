function createConfig() {
    return {
        home: "about/01_overview",
        release: "3.0.0-2006051343",
        releases: [
            "3.0.0-2006051343"
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
            path: '/about/01_overview',
            meta: {
                h1: 'Overview',
                title: 'Overview',
                h1Prefix: null,
                description: 'Coherence Operator documentation',
                keywords: 'oracle coherence, kubernetes, operator, documentation',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('about-01_overview', '/about/01_overview', {})
        },
        {
            path: '/about/02_introduction',
            meta: {
                h1: 'Coherence Operator Introduction',
                title: 'Coherence Operator Introduction',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('about-02_introduction', '/about/02_introduction', {})
        },
        {
            path: '/about/03_quickstart',
            meta: {
                h1: 'Quick Start',
                title: 'Quick Start',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('about-03_quickstart', '/about/03_quickstart', {})
        },
        {
            path: '/about/04_coherence_spec',
            meta: {
                h1: 'Coherence Operator API Docs',
                title: 'Coherence Operator API Docs',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('about-04_coherence_spec', '/about/04_coherence_spec', {})
        },
        {
            path: '/installation/01_installation',
            meta: {
                h1: 'Coherence Operator Installation',
                title: 'Coherence Operator Installation',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('installation-01_installation', '/installation/01_installation', {})
        },
        {
            path: '/installation/02_pre_release_versions',
            meta: {
                h1: 'Accessing Pre-Release Versions',
                title: 'Accessing Pre-Release Versions',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('installation-02_pre_release_versions', '/installation/02_pre_release_versions', {})
        },
        {
            path: '/installation/04_obtain_coherence_images',
            meta: {
                h1: 'Obtain Coherence Images',
                title: 'Obtain Coherence Images',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('installation-04_obtain_coherence_images', '/installation/04_obtain_coherence_images', {})
        },
        {
            path: '/installation/05_private_repos',
            meta: {
                h1: 'Using Private Image Registries',
                title: 'Using Private Image Registries',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('installation-05_private_repos', '/installation/05_private_repos', {})
        },
        {
            path: '/installation/06_openshift',
            meta: {
                h1: 'Coherence Clusters on OpenShift',
                title: 'Coherence Clusters on OpenShift',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('installation-06_openshift', '/installation/06_openshift', {})
        },
        {
            path: '/applications/010_overview',
            meta: {
                h1: 'Overview',
                title: 'Overview',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('applications-010_overview', '/applications/010_overview', {})
        },
        {
            path: '/applications/020_build_application',
            meta: {
                h1: 'Build Application Images',
                title: 'Build Application Images',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('applications-020_build_application', '/applications/020_build_application', {})
        },
        {
            path: '/applications/030_deploy_application',
            meta: {
                h1: 'Deploy Coherence Applications',
                title: 'Deploy Coherence Applications',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('applications-030_deploy_application', '/applications/030_deploy_application', {})
        },
        {
            path: '/applications/040_application_main',
            meta: {
                h1: 'Set the Application Main',
                title: 'Set the Application Main',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('applications-040_application_main', '/applications/040_application_main', {})
        },
        {
            path: '/applications/050_application_args',
            meta: {
                h1: 'Set Application Arguments',
                title: 'Set Application Arguments',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('applications-050_application_args', '/applications/050_application_args', {})
        },
        {
            path: '/applications/060_application_working_dir',
            meta: {
                h1: 'Set the Working Directory',
                title: 'Set the Working Directory',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('applications-060_application_working_dir', '/applications/060_application_working_dir', {})
        },
        {
            path: '/coherence_settings/010_overview',
            meta: {
                h1: 'Configuring Coherence',
                title: 'Configuring Coherence',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('coherence_settings-010_overview', '/coherence_settings/010_overview', {})
        },
        {
            path: '/coherence_settings/020_cluster_name',
            meta: {
                h1: 'Coherence Cluster Name',
                title: 'Coherence Cluster Name',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('coherence_settings-020_cluster_name', '/coherence_settings/020_cluster_name', {})
        },
        {
            path: '/coherence_settings/030_cache_config',
            meta: {
                h1: 'Cache Configuration File',
                title: 'Cache Configuration File',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('coherence_settings-030_cache_config', '/coherence_settings/030_cache_config', {})
        },
        {
            path: '/coherence_settings/040_override_file',
            meta: {
                h1: 'Operational Configuration File',
                title: 'Operational Configuration File',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('coherence_settings-040_override_file', '/coherence_settings/040_override_file', {})
        },
        {
            path: '/coherence_settings/050_storage_enabled',
            meta: {
                h1: 'Storage Enabled or Disabled',
                title: 'Storage Enabled or Disabled',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('coherence_settings-050_storage_enabled', '/coherence_settings/050_storage_enabled', {})
        },
        {
            path: '/coherence_settings/060_log_level',
            meta: {
                h1: 'Coherence Log Level',
                title: 'Coherence Log Level',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('coherence_settings-060_log_level', '/coherence_settings/060_log_level', {})
        },
        {
            path: '/coherence_settings/070_wka',
            meta: {
                h1: 'Well Known Addressing',
                title: 'Well Known Addressing',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('coherence_settings-070_wka', '/coherence_settings/070_wka', {})
        },
        {
            path: '/coherence_settings/080_persistence',
            meta: {
                h1: 'Coherence Persistence',
                title: 'Coherence Persistence',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('coherence_settings-080_persistence', '/coherence_settings/080_persistence', {})
        },
        {
            path: '/jvm_settings/010_overview',
            meta: {
                h1: 'Overview',
                title: 'Overview',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('jvm_settings-010_overview', '/jvm_settings/010_overview', {})
        },
        {
            path: '/jvm_settings/020_classpath',
            meta: {
                h1: 'Set the Classpath',
                title: 'Set the Classpath',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('jvm_settings-020_classpath', '/jvm_settings/020_classpath', {})
        },
        {
            path: '/jvm_settings/030_jvm_args',
            meta: {
                h1: 'Arbitrary JVM Arguments',
                title: 'Arbitrary JVM Arguments',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('jvm_settings-030_jvm_args', '/jvm_settings/030_jvm_args', {})
        },
        {
            path: '/jvm_settings/040_gc',
            meta: {
                h1: 'Garbage Collector Settings',
                title: 'Garbage Collector Settings',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('jvm_settings-040_gc', '/jvm_settings/040_gc', {})
        },
        {
            path: '/jvm_settings/050_memory',
            meta: {
                h1: 'Heap & Memory Settings',
                title: 'Heap & Memory Settings',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('jvm_settings-050_memory', '/jvm_settings/050_memory', {})
        },
        {
            path: '/jvm_settings/070_debugger',
            meta: {
                h1: 'Debugger Configuration',
                title: 'Debugger Configuration',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('jvm_settings-070_debugger', '/jvm_settings/070_debugger', {})
        },
        {
            path: '/jvm_settings/080_jmx',
            meta: {
                h1: 'Using JMX',
                title: 'Using JMX',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('jvm_settings-080_jmx', '/jvm_settings/080_jmx', {})
        },
        {
            path: '/jvm_settings/090_container_limits',
            meta: {
                h1: 'Container Resource Limits',
                title: 'Container Resource Limits',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('jvm_settings-090_container_limits', '/jvm_settings/090_container_limits', {})
        },
        {
            path: '/ports/010_overview',
            meta: {
                h1: 'Overview',
                title: 'Overview',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('ports-010_overview', '/ports/010_overview', {})
        },
        {
            path: '/ports/020_container_ports',
            meta: {
                h1: 'Additional Container Ports',
                title: 'Additional Container Ports',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('ports-020_container_ports', '/ports/020_container_ports', {})
        },
        {
            path: '/ports/030_services',
            meta: {
                h1: 'Configure Services for Ports',
                title: 'Configure Services for Ports',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('ports-030_services', '/ports/030_services', {})
        },
        {
            path: '/ports/040_servicemonitors',
            meta: {
                h1: 'Prometheus ServiceMonitors',
                title: 'Prometheus ServiceMonitors',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('ports-040_servicemonitors', '/ports/040_servicemonitors', {})
        },
        {
            path: '/scaling/010_overview',
            meta: {
                h1: 'Scale Coherence Deployments',
                title: 'Scale Coherence Deployments',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('scaling-010_overview', '/scaling/010_overview', {})
        },
        {
            path: '/ordering/010_overview',
            meta: {
                h1: 'Deployment Start Order',
                title: 'Deployment Start Order',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('ordering-010_overview', '/ordering/010_overview', {})
        },
        {
            path: '/management_and_diagnostics/010_overview',
            meta: {
                h1: 'Overview',
                title: 'Overview',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('management_and_diagnostics-010_overview', '/management_and_diagnostics/010_overview', {})
        },
        {
            path: '/management_and_diagnostics/020_management_over_rest',
            meta: {
                h1: 'Management over REST',
                title: 'Management over REST',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('management_and_diagnostics-020_management_over_rest', '/management_and_diagnostics/020_management_over_rest', {})
        },
        {
            path: '/management_and_diagnostics/030_visualvm',
            meta: {
                h1: 'Using VisualVM',
                title: 'Using VisualVM',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('management_and_diagnostics-030_visualvm', '/management_and_diagnostics/030_visualvm', {})
        },
        {
            path: '/management_and_diagnostics/040_ssl',
            meta: {
                h1: 'SSL with Management over REST',
                title: 'SSL with Management over REST',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('management_and_diagnostics-040_ssl', '/management_and_diagnostics/040_ssl', {})
        },
        {
            path: '/metrics/010_overview',
            meta: {
                h1: 'Overview',
                title: 'Overview',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('metrics-010_overview', '/metrics/010_overview', {})
        },
        {
            path: '/metrics/020_metrics',
            meta: {
                h1: 'Publish Metrics',
                title: 'Publish Metrics',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('metrics-020_metrics', '/metrics/020_metrics', {})
        },
        {
            path: '/metrics/030_importing',
            meta: {
                h1: 'Importing Grafana Dashboards',
                title: 'Importing Grafana Dashboards',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('metrics-030_importing', '/metrics/030_importing', {})
        },
        {
            path: '/metrics/040_dashboards',
            meta: {
                h1: 'Grafana Dashboards',
                title: 'Grafana Dashboards',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('metrics-040_dashboards', '/metrics/040_dashboards', {})
        },
        {
            path: '/metrics/050_ssl',
            meta: {
                h1: 'SSL with Metrics',
                title: 'SSL with Metrics',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('metrics-050_ssl', '/metrics/050_ssl', {})
        },
        {
            path: '/logging/010_overview',
            meta: {
                h1: 'Overview',
                title: 'Overview',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('logging-010_overview', '/logging/010_overview', {})
        },
        {
            path: '/logging/020_logging',
            meta: {
                h1: 'Log Capture with Fluentd',
                title: 'Log Capture with Fluentd',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('logging-020_logging', '/logging/020_logging', {})
        },
        {
            path: '/logging/030_kibana',
            meta: {
                h1: 'Using Kibana Dashboards',
                title: 'Using Kibana Dashboards',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('logging-030_kibana', '/logging/030_kibana', {})
        },
        {
            path: '/misc_pod_settings/010_overview',
            meta: {
                h1: 'Overview',
                title: 'Overview',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('misc_pod_settings-010_overview', '/misc_pod_settings/010_overview', {})
        },
        {
            path: '/misc_pod_settings/020_environment',
            meta: {
                h1: 'Environment Variables',
                title: 'Environment Variables',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('misc_pod_settings-020_environment', '/misc_pod_settings/020_environment', {})
        },
        {
            path: '/misc_pod_settings/030_labels',
            meta: {
                h1: 'Pod Labels',
                title: 'Pod Labels',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('misc_pod_settings-030_labels', '/misc_pod_settings/030_labels', {})
        },
        {
            path: '/misc_pod_settings/040_annotations',
            meta: {
                h1: 'Pod Annotations',
                title: 'Pod Annotations',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('misc_pod_settings-040_annotations', '/misc_pod_settings/040_annotations', {})
        },
        {
            path: '/misc_pod_settings/050_configmap_volumes',
            meta: {
                h1: 'Add ConfigMap Volumes',
                title: 'Add ConfigMap Volumes',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('misc_pod_settings-050_configmap_volumes', '/misc_pod_settings/050_configmap_volumes', {})
        },
        {
            path: '/misc_pod_settings/060_secret_volumes',
            meta: {
                h1: 'Add Secrets Volumes',
                title: 'Add Secrets Volumes',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('misc_pod_settings-060_secret_volumes', '/misc_pod_settings/060_secret_volumes', {})
        },
        {
            path: '/misc_pod_settings/070_add_volumes',
            meta: {
                h1: 'Add Pod Volumes',
                title: 'Add Pod Volumes',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('misc_pod_settings-070_add_volumes', '/misc_pod_settings/070_add_volumes', {})
        },
        {
            path: '/misc_pod_settings/080_add_containers',
            meta: {
                h1: 'Configure Additional Containers',
                title: 'Configure Additional Containers',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('misc_pod_settings-080_add_containers', '/misc_pod_settings/080_add_containers', {})
        },
        {
            path: '/misc_pod_settings/090_pod_scheduling',
            meta: {
                h1: 'Configure Pod Scheduling',
                title: 'Configure Pod Scheduling',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('misc_pod_settings-090_pod_scheduling', '/misc_pod_settings/090_pod_scheduling', {})
        },
        {
            path: '/misc_pod_settings/100_resources',
            meta: {
                h1: 'Container Resource Limits',
                title: 'Container Resource Limits',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('misc_pod_settings-100_resources', '/misc_pod_settings/100_resources', {})
        },
        {
            path: '/misc_pod_settings/110_readiness',
            meta: {
                h1: 'Readiness & Liveness Probes',
                title: 'Readiness & Liveness Probes',
                h1Prefix: null,
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('misc_pod_settings-110_readiness', '/misc_pod_settings/110_readiness', {})
        },
        {
            path: '/', redirect: '/about/01_overview'
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
                { href: '/about/01_overview', title: 'Overview' },
                { href: '/about/02_introduction', title: 'Coherence Operator Introduction' },
                { href: '/about/03_quickstart', title: 'Quick Start' },
                { href: '/about/04_coherence_spec', title: 'Coherence Operator API Docs' }
            ]
        },
        {
            title: 'Installation',
            action: 'fa-save',
            group: '/install',
            items: [
                { href: '/installation/01_installation', title: 'Coherence Operator Installation' },
                { href: '/installation/02_pre_release_versions', title: 'Accessing Pre-Release Versions' },
                { href: '/installation/04_obtain_coherence_images', title: 'Obtain Coherence Images' },
                { href: '/installation/05_private_repos', title: 'Using Private Image Registries' },
                { href: '/installation/06_openshift', title: 'Coherence Clusters on OpenShift' }
            ]
        },
        {
            title: 'Deploy Applications',
            action: 'cloud_upload',
            group: '/applications',
            items: [
                { href: '/applications/010_overview', title: 'Overview' },
                { href: '/applications/020_build_application', title: 'Build Application Images' },
                { href: '/applications/030_deploy_application', title: 'Deploy Coherence Applications' },
                { href: '/applications/040_application_main', title: 'Set the Application Main' },
                { href: '/applications/050_application_args', title: 'Set Application Arguments' },
                { href: '/applications/060_application_working_dir', title: 'Set the Working Directory' }
            ]
        },
        {
            title: 'Coherence Settings',
            action: 'fa-cogs',
            group: '/coherence',
            items: [
                { href: '/coherence_settings/010_overview', title: 'Configuring Coherence' },
                { href: '/coherence_settings/020_cluster_name', title: 'Coherence Cluster Name' },
                { href: '/coherence_settings/030_cache_config', title: 'Cache Configuration File' },
                { href: '/coherence_settings/040_override_file', title: 'Operational Configuration File' },
                { href: '/coherence_settings/050_storage_enabled', title: 'Storage Enabled or Disabled' },
                { href: '/coherence_settings/060_log_level', title: 'Coherence Log Level' },
                { href: '/coherence_settings/070_wka', title: 'Well Known Addressing' },
                { href: '/coherence_settings/080_persistence', title: 'Coherence Persistence' }
            ]
        },
        {
            title: 'JVM Settings',
            action: 'fa-cog',
            group: '/jvm',
            items: [
                { href: '/jvm_settings/010_overview', title: 'Overview' },
                { href: '/jvm_settings/020_classpath', title: 'Set the Classpath' },
                { href: '/jvm_settings/030_jvm_args', title: 'Arbitrary JVM Arguments' },
                { href: '/jvm_settings/040_gc', title: 'Garbage Collector Settings' },
                { href: '/jvm_settings/050_memory', title: 'Heap & Memory Settings' },
                { href: '/jvm_settings/070_debugger', title: 'Debugger Configuration' },
                { href: '/jvm_settings/080_jmx', title: 'Using JMX' },
                { href: '/jvm_settings/090_container_limits', title: 'Container Resource Limits' }
            ]
        },
        {
            title: 'Expose Ports & Services',
            action: 'control_camera',
            group: '/ports',
            items: [
                { href: '/ports/010_overview', title: 'Overview' },
                { href: '/ports/020_container_ports', title: 'Additional Container Ports' },
                { href: '/ports/030_services', title: 'Configure Services for Ports' },
                { href: '/ports/040_servicemonitors', title: 'Prometheus ServiceMonitors' }
            ]
        },
        {
            title: 'Scaling Up & Down',
            action: 'fa-balance-scale',
            group: '/scaling',
            items: [
                { href: '/scaling/010_overview', title: 'Scale Coherence Deployments' }
            ]
        },
        {
            title: 'Start-up Order',
            action: 'line_weight',
            group: '/ordering',
            items: [
                { href: '/ordering/010_overview', title: 'Deployment Start Order' }
            ]
        },
        {
            title: 'Management Diagnostics',
            action: 'fa-stethoscope',
            group: '/management',
            items: [
                { href: '/management_and_diagnostics/010_overview', title: 'Overview' },
                { href: '/management_and_diagnostics/020_management_over_rest', title: 'Management over REST' },
                { href: '/management_and_diagnostics/030_visualvm', title: 'Using VisualVM' },
                { href: '/management_and_diagnostics/040_ssl', title: 'SSL with Management over REST' }
            ]
        },
        {
            title: 'Metrics',
            action: 'speed',
            group: '/metrics',
            items: [
                { href: '/metrics/010_overview', title: 'Overview' },
                { href: '/metrics/020_metrics', title: 'Publish Metrics' },
                { href: '/metrics/030_importing', title: 'Importing Grafana Dashboards' },
                { href: '/metrics/040_dashboards', title: 'Grafana Dashboards' },
                { href: '/metrics/050_ssl', title: 'SSL with Metrics' }
            ]
        },
        {
            title: 'Logging',
            action: 'find_in_page',
            group: '/logging',
            items: [
                { href: '/logging/010_overview', title: 'Overview' },
                { href: '/logging/020_logging', title: 'Log Capture with Fluentd' },
                { href: '/logging/030_kibana', title: 'Using Kibana Dashboards' }
            ]
        },
        {
            title: 'Other Pod Settings',
            action: 'widgets',
            group: '/misc',
            items: [
                { href: '/misc_pod_settings/010_overview', title: 'Overview' },
                { href: '/misc_pod_settings/020_environment', title: 'Environment Variables' },
                { href: '/misc_pod_settings/030_labels', title: 'Pod Labels' },
                { href: '/misc_pod_settings/040_annotations', title: 'Pod Annotations' },
                { href: '/misc_pod_settings/050_configmap_volumes', title: 'Add ConfigMap Volumes' },
                { href: '/misc_pod_settings/060_secret_volumes', title: 'Add Secrets Volumes' },
                { href: '/misc_pod_settings/070_add_volumes', title: 'Add Pod Volumes' },
                { href: '/misc_pod_settings/080_add_containers', title: 'Configure Additional Containers' },
                { href: '/misc_pod_settings/090_pod_scheduling', title: 'Configure Pod Scheduling' },
                { href: '/misc_pod_settings/100_resources', title: 'Container Resource Limits' },
                { href: '/misc_pod_settings/110_readiness', title: 'Readiness & Liveness Probes' }
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
            href: 'https://coherence.java.net',
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