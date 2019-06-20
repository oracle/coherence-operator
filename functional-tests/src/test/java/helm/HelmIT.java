/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.yaml.YAMLFactory;
import com.oracle.bedrock.runtime.Application;
import com.oracle.bedrock.runtime.console.CapturingApplicationConsole;
import com.oracle.bedrock.runtime.k8s.helm.HelmInstall;
import com.oracle.bedrock.runtime.options.Console;
import com.tangosol.util.Resources;
import org.junit.BeforeClass;
import org.junit.ClassRule;
import org.junit.Test;
import org.junit.rules.TemporaryFolder;

import java.io.File;
import java.net.URL;
import java.util.Arrays;
import java.util.Collections;
import java.util.HashMap;
import java.util.Hashtable;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.Queue;
import java.util.UUID;

import static helm.HelmUtils.HELM_TIMEOUT;
import static org.hamcrest.CoreMatchers.containsString;
import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.CoreMatchers.notNullValue;
import static org.hamcrest.CoreMatchers.nullValue;
import static org.junit.Assert.assertThat;
import static org.junit.Assert.fail;


/**
 * This test performs helm install --debug --dry-run commands and verifies the
 * generated yaml for various values settings.
 *
 * @author jk
 */
public class HelmIT
        extends BaseHelmChartTest
    {
    /**
     * Unpack the Coherence chart into a temporary folder.
     *
     * @throws Exception  if the chart cannot be extracted
     */
    @BeforeClass
    public static void setup() throws Exception
        {
        s_fileChart = extractChart(COHERENCE_HELM_CHART_NAME, COHERENCE_HELM_CHART_URL);

        // install the chart without any values to obtain the yaml for a default install
        s_mapDefaultYaml = install("default-release");
        }

    // ----- test methods ---------------------------------------------------

    @Test
    public void shouldAllowZeroPullSecrets()
        {
        LinkedHashMap mapCoh = s_mapDefaultYaml.get("coherence.yaml");

        assertThat(mapCoh, is(notNullValue()));

        Object oResult = getYamlValue(mapCoh, "spec", "template", "spec", "imagePullSecrets");

        assertThat(oResult, is(nullValue()));
        }

    @Test
    public void shouldAllowSinglePullSecretAsString()
        {
        String                     sRelease = "test-release";
        Map<String, LinkedHashMap> map      = install(sRelease, "imagePullSecrets=foo");
        LinkedHashMap              mapCoh   = map.get("coherence.yaml");

        assertThat(mapCoh, is(notNullValue()));

        Object oResult = getYamlValue(mapCoh, "spec", "template", "spec", "imagePullSecrets");

        Map<String, String> mapFoo = Collections.singletonMap("name", "foo");
        List                list   = Arrays.asList(mapFoo);

        assertThat(oResult, is(list));
        }

    @Test
    public void shouldAllowSinglePullSecretAsList()
        {
        String                     sRelease = "test-release";
        Map<String, LinkedHashMap> map      = install(sRelease, "imagePullSecrets={foo}");
        LinkedHashMap              mapCoh   = map.get("coherence.yaml");

        assertThat(mapCoh, is(notNullValue()));

        Object oResult = getYamlValue(mapCoh, "spec", "template", "spec", "imagePullSecrets");

        Map<String, String> mapFoo = Collections.singletonMap("name", "foo");
        List                list   = Arrays.asList(mapFoo);

        assertThat(oResult, is(list));
        }

    @Test
    public void shouldAllowMultiplePullSecrets()
        {
        String                     sRelease = "test-release";
        Map<String, LinkedHashMap> map      = install(sRelease, "imagePullSecrets={foo,bar}");
        LinkedHashMap              mapCoh   = map.get("coherence.yaml");

        assertThat(mapCoh, is(notNullValue()));

        Object oResult = getYamlValue(mapCoh, "spec", "template", "spec", "imagePullSecrets");

        Map<String, String> mapFoo = Collections.singletonMap("name", "foo");
        Map<String, String> mapBar = Collections.singletonMap("name", "bar");
        List                list   = Arrays.asList(mapFoo, mapBar);

        assertThat(oResult, is(list));
        }

    @Test
    public void shouldSetClusterNameToSpecifiedName()
        {
        String                     sRelease = "test-release";
        String                     sCluster = "test-cluster";
        Map<String, LinkedHashMap> map      = install(sRelease, "cluster=" + sCluster);

        assertEnvironmentVariable(map, "COH_CLUSTER_NAME", sCluster);
        }

    @Test
    public void shouldSetClusterNameToReleaseNameByDefault()
        {
        assertEnvironmentVariable(s_mapDefaultYaml, "COH_CLUSTER_NAME", "default-release");
        }

    @Test
    public void shouldNotSetCacheConfigurationByDefault()
        {
        assertEnvironmentVariableNotSet(s_mapDefaultYaml, "COH_CACHE_CONFIG");
        }

    @Test
    public void shouldSetCacheConfiguration()
        {
        String                     sRelease = "test-release";
        String                     sConfig  = "test-config.xml";
        Map<String, LinkedHashMap> map      = install(sRelease, "store.cacheConfig=" + sConfig);

        assertEnvironmentVariable(map, "COH_CACHE_CONFIG", sConfig);
        }

    @Test
    public void shouldEnablePofAndPofConfig()
        {
        String                     sRelease = "test-release";
        String                     sConfig  = "test-pof-config.xml";
        Map<String, LinkedHashMap> map      = install(sRelease, "store.pof.config=" + sConfig);

        assertEnvironmentVariable(map, "COH_POF_CONFIG", sConfig);
        }

    @Test
    public void shouldEnablePofWithDefaultPofConfig()
        {
        String                     sRelease = "test-release";
        Map<String, LinkedHashMap> map      = install(sRelease);

        assertEnvironmentVariable(map, "COH_POF_CONFIG", "pof-config.xml");
        }

    @Test
    public void shouldNotSetLogLevelByDefault()
        {
        assertEnvironmentVariableNotSet(s_mapDefaultYaml, "COH_LOG_LEVEL");
        }

    @Test
    public void shouldSetLogLevel()
        {
        String                     sRelease = "test-release";
        Map<String, LinkedHashMap> map      = install(sRelease, "store.logging.level=9");

        assertEnvironmentVariable(map, "COH_LOG_LEVEL", "9");
        }

    @Test
    public void shouldUseDefaultLoggingConfiguration()
        {
        assertEnvironmentVariable(s_mapDefaultYaml, "COH_LOGGING_CONFIG", "/scripts/logging.properties");
        }

    @Test
    public void shouldUseConfigMapLoggingConfiguration()
        {
        String                     sRelease = "test-release";
        Map<String, LinkedHashMap> map      = install(sRelease, "store.logging.configFile=foo.properties", "store.logging.configMapName=FooMap");
        LinkedHashMap              mapCoh     = map.get("coherence.yaml");

        assertThat(mapCoh, is(notNullValue()));

        assertEnvironmentVariable(map, "COH_LOGGING_CONFIG", "/loggingconfig/foo.properties");

        List   listContainer  = getYamlValue(mapCoh, "spec", "template", "spec", "containers");
        Map    mapContainer   = getNamedYamlValue(listContainer, "coherence");
        List   listMounts     = getYamlValue(mapContainer, "volumeMounts");
        Map    mapVolumeMount = getNamedYamlValue(listMounts, "logging-config");

        assertThat("Should have added logging-config volume mount", mapVolumeMount, is(notNullValue()));

        List listVolumes = getYamlValue(mapCoh, "spec", "template", "spec", "volumes");
        Map  mapVolume   = getNamedYamlValue(listVolumes, "logging-config");

        assertThat("Should have added logging-config volume", mapVolume, is(notNullValue()));
        }

    @Test
    public void shouldUseUserArtifactsLoggingConfiguration()
        {
        String                     sRelease = "test-release";
        Map<String, LinkedHashMap> map      = install(sRelease, "store.logging.configFile=foo.properties", "userArtifacts.image=foo/image:1.0");
        LinkedHashMap              mapCoh     = map.get("coherence.yaml");

        assertThat(mapCoh, is(notNullValue()));

        assertEnvironmentVariable(map, "COH_LOGGING_CONFIG", "/u01/oracle/oracle_home/coherence/ext/conf/foo.properties");
        }

    @Test
    public void shouldUseLoggingConfigurationValue()
        {
        String                     sRelease = "test-release";
        Map<String, LinkedHashMap> map      = install(sRelease, "store.logging.configFile=/home/root/foo.properties");
        LinkedHashMap              mapCoh     = map.get("coherence.yaml");

        assertThat(mapCoh, is(notNullValue()));

        assertEnvironmentVariable(map, "COH_LOGGING_CONFIG", "/home/root/foo.properties");
        }

    @Test
    public void shouldNotConfigureLocalStorageByDefault()
        {
        assertEnvironmentVariableNotSet(s_mapDefaultYaml, "COH_STORAGE_ENABLED");
        }

    @Test
    public void shouldNotConfigureLocalStorageIfValueIsNull()
        {
        String                     sRelease = "test-release";
        Map<String, LinkedHashMap> map      = install(sRelease, "store.storageEnabled=null");

        assertEnvironmentVariableNotSet(map, "COH_STORAGE_ENABLED");
        }

    @Test
    public void shouldConfigureLocalStorageToTrue()
        {
        String                     sRelease = "test-release";
        Map<String, LinkedHashMap> map      = install(sRelease, "store.storageEnabled=true");

        assertEnvironmentVariable(map, "COH_STORAGE_ENABLED", "true");
        }

    @Test
    public void shouldConfigureLocalStorageToFalse()
        {
        String                     sRelease = "test-release";
        Map<String, LinkedHashMap> map      = install(sRelease, "store.storageEnabled=false");

        assertEnvironmentVariable(map, "COH_STORAGE_ENABLED", "false");
        }

    @Test
    public void shouldNotSetServiceAccountNameByDefault()
        {
        LinkedHashMap mapCoh = s_mapDefaultYaml.get("coherence.yaml");

        assertThat(mapCoh, is(notNullValue()));

        String sName = getYamlValue(mapCoh, "spec", "template", "spec", "serviceAccountName");

        assertThat(sName, is(nullValue()));
        }

    @Test
    public void shouldSetServiceAccountName()
        {
        String                     sRelease = "test-release";
        Map<String, LinkedHashMap> map      = install(sRelease, "serviceAccountName=foo");
        LinkedHashMap              mapCoh   = map.get("coherence.yaml");

        assertThat(mapCoh, is(notNullValue()));

        String sName = getYamlValue(mapCoh, "spec", "template", "spec", "serviceAccountName");

        assertThat(sName, is("foo"));
        }

    @Test
    public void shouldNotSetServiceAccountNameToDefaultAccount()
        {
        String                     sRelease = "test-release";
        Map<String, LinkedHashMap> map      = install(sRelease, "serviceAccountName=default");
        LinkedHashMap              mapCoh   = map.get("coherence.yaml");

        assertThat(mapCoh, is(notNullValue()));

        String sName = getYamlValue(mapCoh, "spec", "template", "spec", "serviceAccountName");

        assertThat(sName, is(nullValue()));
        }

    @Test
    public void shouldAddSingleExtraPort() throws Exception
        {
        URL                        url        = Resources.findFileOrResource("values/helm-values-with-single-extra-port.yaml", null);
        String                     sRelease   = "test-release";
        Map<String, LinkedHashMap> map        = install(sRelease, new File(url.toURI()));
        LinkedHashMap              mapCoh     = map.get("coherence.yaml");
        LinkedHashMap              mapService = map.get("service.yaml");

        assertThat(mapCoh, is(notNullValue()));

        List   listContainer = getYamlValue(mapCoh, "spec", "template", "spec", "containers");
        Map    mapContainer  = getNamedYamlValue(listContainer, "coherence");
        List   listPorts     = getYamlValue(mapContainer, "ports");
        Map    mapFoo        = getNamedYamlValue(listPorts, "foo");

        assertThat(mapFoo, is(notNullValue()));
        assertThat(mapFoo.get("containerPort"), is(8080));

        List listServicePorts = getYamlValue(mapService, "spec", "ports");
        mapFoo = getNamedYamlValue(listServicePorts, "foo");

        assertThat(mapFoo, is(notNullValue()));
        assertThat(mapFoo.get("port"), is(8080));
        assertThat(mapFoo.get("protocol"), is("TCP"));
        assertThat(mapFoo.get("targetPort"), is("foo"));
        }

    @Test
    public void shouldAddMultipleExtraPorts() throws Exception
        {
        URL                        url        = Resources.findFileOrResource("values/helm-values-with-extra-ports.yaml", null);
        String                     sRelease   = "test-release";
        Map<String, LinkedHashMap> map        = install(sRelease, new File(url.toURI()));
        LinkedHashMap              mapCoh     = map.get("coherence.yaml");
        LinkedHashMap              mapService = map.get("service.yaml");

        assertThat(mapCoh, is(notNullValue()));
        assertThat(mapService, is(notNullValue()));

        List   listContainer = getYamlValue(mapCoh, "spec", "template", "spec", "containers");
        Map    mapContainer  = getNamedYamlValue(listContainer, "coherence");
        List   listPorts     = getYamlValue(mapContainer, "ports");
        Map    mapFoo        = getNamedYamlValue(listPorts, "foo");
        Map    mapBar        = getNamedYamlValue(listPorts, "bar");

        assertThat(mapFoo, is(notNullValue()));
        assertThat(mapFoo.get("containerPort"), is(8080));
        assertThat(mapBar, is(notNullValue()));
        assertThat(mapBar.get("containerPort"), is(1234));

        List listServicePorts = getYamlValue(mapService, "spec", "ports");
        mapFoo = getNamedYamlValue(listServicePorts, "foo");
        mapBar = getNamedYamlValue(listServicePorts, "bar");

        assertThat(mapFoo, is(notNullValue()));
        assertThat(mapFoo.get("port"), is(8080));
        assertThat(mapFoo.get("protocol"), is("TCP"));
        assertThat(mapFoo.get("targetPort"), is("foo"));

        assertThat(mapBar, is(notNullValue()));
        assertThat(mapBar.get("port"), is(1234));
        assertThat(mapBar.get("protocol"), is("TCP"));
        assertThat(mapBar.get("targetPort"), is("bar"));
        }


    @Test
    @SuppressWarnings("unchecked")
    public void shouldNotSetSetExtraPVCFieldsByDefault()
        {
        String                     sRelease = "test-release";
        Map<String, LinkedHashMap> map      = install(sRelease, "store.persistence.enabled=true", "store.snapshot.enabled=true");
        LinkedHashMap              mapCoh   = map.get("coherence.yaml");

        assertThat(mapCoh, is(notNullValue()));

        List listPVC = getYamlValue(mapCoh, "spec", "volumeClaimTemplates");

        assertThat(listPVC, is(notNullValue()));

        Map  mapPersistence = getPVC(listPVC, "persistence-volume");
        Map  mapSnapshot    = getPVC(listPVC, "snapshot-volume");

        assertThat(mapPersistence, is(notNullValue()));
        assertThat(mapSnapshot, is(notNullValue()));

        Map mapSpecPersistence = (Map) mapPersistence.get("spec");
        Map mapSpecSnapshot    = (Map) mapPersistence.get("spec");

        Map mapExpectedPersistence = new HashMap();

        mapExpectedPersistence.put("accessModes", Collections.singletonList("ReadWriteOnce"));
        mapExpectedPersistence.put("resources", Collections.singletonMap("requests", Collections.singletonMap("storage", "2Gi")));

        Map mapExpectedSnapshot = new HashMap();

        mapExpectedSnapshot.put("accessModes", Collections.singletonList("ReadWriteOnce"));
        mapExpectedSnapshot.put("resources", Collections.singletonMap("requests", Collections.singletonMap("storage", "2Gi")));

        assertThat(mapSpecPersistence, is(mapExpectedPersistence));
        assertThat(mapSpecSnapshot, is(mapExpectedSnapshot));
        }

    @Test
    @SuppressWarnings("unchecked")
    public void shouldSetSetExtraPVCFields() throws Exception
        {
        String                     sRelease = "test-release";
        URL                        url      = Resources.findFileOrResource("values/helm-values-pvc-fields.yaml", null);
        Map<String, LinkedHashMap> map      = install(sRelease, new File(url.toURI()));
        LinkedHashMap              mapCoh   = map.get("coherence.yaml");

        assertThat(mapCoh, is(notNullValue()));

        List listPVC = getYamlValue(mapCoh, "spec", "volumeClaimTemplates");

        assertThat(listPVC, is(notNullValue()));

        Map  mapPersistence = getPVC(listPVC, "persistence-volume");
        Map  mapSnapshot    = getPVC(listPVC, "snapshot-volume");

        assertThat(mapPersistence, is(notNullValue()));
        assertThat(mapSnapshot, is(notNullValue()));

        Map mapSpecPersistence = (Map) mapPersistence.get("spec");
        Map mapSpecSnapshot    = (Map) mapSnapshot.get("spec");

        Map mapExpectedPersistence = new LinkedHashMap();
        Map mapSelectorPersistence = new LinkedHashMap();

        mapExpectedPersistence.put("accessModes", Collections.singletonList("ReadWriteOnce"));
        mapExpectedPersistence.put("storageClassName", "test-class-one");
        mapExpectedPersistence.put("dataSource", "test-data-source-one");
        mapExpectedPersistence.put("volumeMode", "test-volume-mode-one");
        mapExpectedPersistence.put("volumeName", "test-volume-name-one");
        mapSelectorPersistence.put("selector-one", "label-one");
        mapSelectorPersistence.put("selector-two", "label-two");
        mapExpectedPersistence.put("selector", Collections.singletonMap("matchLabels", mapSelectorPersistence));
        mapExpectedPersistence.put("resources", Collections.singletonMap("requests", Collections.singletonMap("storage", "2Gi")));

        Map mapExpectedSnapshot = new LinkedHashMap();
        Map mapSelectorSnapshot = new LinkedHashMap();

        mapExpectedSnapshot.put("accessModes", Collections.singletonList("ReadWriteOnce"));
        mapExpectedSnapshot.put("storageClassName", "test-class-two");
        mapExpectedSnapshot.put("dataSource", "test-data-source-two");
        mapExpectedSnapshot.put("volumeMode", "test-volume-mode-two");
        mapExpectedSnapshot.put("volumeName", "test-volume-name-two");
        mapSelectorSnapshot.put("selector-three", "label-three");
        mapSelectorSnapshot.put("selector-four", "label-four");
        mapExpectedSnapshot.put("selector", Collections.singletonMap("matchLabels", mapSelectorSnapshot));
        mapExpectedSnapshot.put("resources", Collections.singletonMap("requests", Collections.singletonMap("storage", "2Gi")));

        assertThat(mapSpecPersistence, is(mapExpectedPersistence));
        assertThat(mapSpecSnapshot, is(mapExpectedSnapshot));
        }

    // ----- helper methods -------------------------------------------------

    @SuppressWarnings("unchecked")
    private Map getPVC(List list, String sName)
        {
        for (Map mapPVC : (List<Map>)list)
            {
            Map mapMeta = (Map) mapPVC.get("metadata");

            if (mapMeta != null)
                {
                String sNameActual = (String) mapMeta.get("name");

                if (Objects.equals(sName, sNameActual))
                    {
                    return mapPVC;
                    }
                }
            }

        return null;
        }


    private void assertEnvironmentVariable(Map<String, LinkedHashMap> mapYaml,  String sName, String sExpected)
        {
        LinkedHashMap mapCoh = mapYaml.get("coherence.yaml");

        assertThat(mapCoh, is(notNullValue()));

        List   listContainer = getYamlValue(mapCoh, "spec", "template", "spec", "containers");
        Map    mapContainer  = getNamedYamlValue(listContainer, "coherence");
        List   listEnv       = getYamlValue(mapContainer, "env");
        Map    mapOpts       = getNamedYamlValue(listEnv, sName);

        assertThat(mapOpts, is(notNullValue()));

        String sResult = (String) mapOpts.get("value");

        assertThat(sResult, containsString(sExpected));
        }

    private void assertEnvironmentVariableNotSet(Map<String, LinkedHashMap> mapYaml,  String sName)
        {
        LinkedHashMap mapCoh = mapYaml.get("coherence.yaml");

        assertThat(mapCoh, is(notNullValue()));

        List   listContainer = getYamlValue(mapCoh, "spec", "template", "spec", "containers");
        Map    mapContainer  = getNamedYamlValue(listContainer, "coherence");
        List   listEnv       = getYamlValue(mapContainer, "env");
        Map    mapOpts       = getNamedYamlValue(listEnv, sName);
        String sResult       = mapOpts == null ? null : (String) mapOpts.get("value");

        assertThat(sResult, is(nullValue()));
        }

    /**
     * Execute a helm install --debug --dry-run and capture the generated yaml files.
     *
     * @param sRelease  the name to use for the helm release
     *
     * @return  a map of the generated yaml files keyed by template file name
     */
    private static Map<String, LinkedHashMap> install(String sRelease)
        {
        return install(sRelease, new File[0], new String[0]);
        }

    /**
     * Execute a helm install --debug --dry-run and capture the generated yaml files.
     *
     * @param sRelease  the name to use for the helm release
     * @param asValues  the array of values to set
     *
     * @return  a map of the generated yaml files keyed by template file name
     */
    private static Map<String, LinkedHashMap> install(String sRelease, String... asValues)
        {
        return install(sRelease, new File[0], asValues);
        }

    /**
     * Execute a helm install --debug --dry-run and capture the generated yaml files.
     *
     * @param sRelease  the name to use for the helm release
     * @param asValues  the array of values files to use
     *
     * @return  a map of the generated yaml files keyed by template file name
     */
    private static Map<String, LinkedHashMap> install(String sRelease, File... asValues)
        {
        return install(sRelease, asValues, new String[0]);
        }

    /**
     * Execute a helm install --debug --dry-run and capture the generated yaml files.
     *
     * @param sRelease     the name to use for the helm release
     * @param aFileValues  the array of values files to use
     * @param asValues     the array of values to set
     *
     * @return  a map of the generated yaml files keyed by template file name
     */
    private static Map<String, LinkedHashMap> install(String sRelease, File[] aFileValues, String[] asValues)
        {
        try
            {
            CapturingApplicationConsole console = new CapturingApplicationConsole();

            HelmInstall install = s_helm.install(s_fileChart, COHERENCE_HELM_CHART_NAME)
                                        .dryRun()
                                        .debug()
                                        .timeout(HELM_TIMEOUT)
                                        .name(sRelease)
                                        .set("coherenceK8sOperatorTesting=true");

            if (aFileValues != null)
                {
                for (File file : aFileValues)
                    {
                    install = install.values(file);
                    }
                }

            if (asValues != null && asValues.length > 0)
                {
                install = install.set(asValues);
                }

            int nExitCode = -1;
            try (Application app = install.execute(Console.of(console)))
                {
                nExitCode = app.waitFor();
                Thread.sleep(500);  // this sleep is here because there seems to be a delay in getting all of the output
                }

            if (nExitCode != 0)
                {
                int maxRetries = Integer.parseInt(HELM_INSTALL_MAX_RETRY);
                for (int i = maxRetries; nExitCode != 0 && i > 0 ; i--)
                    {
                    System.err.println("Helm install (dry-run) failed with exit code " + nExitCode + " - will retry. "
                        + i + " attempts remaining");

                    logInstallFailure(install, nExitCode, console);

                    console = new CapturingApplicationConsole();
                    try (Application app = install.execute(Console.of(console)))
                        {
                        nExitCode = app.waitFor();
                        Thread.sleep(500);  // this sleep is here because there seems to be a delay in getting all of the output
                        }
                    }
                }

            if (nExitCode != 0)
                {
                HelmUtils.logConsoleOutput("helm-install", console);
                String reason = "Install dry-run failed for helm chart " + COHERENCE_HELM_CHART_NAME
                    + " failed with exit code " + nExitCode;
                fail(reason);
                }

            Queue<String> queue = console.getCapturedOutputLines();
            String        sLine = queue.poll();

            while(sLine != null && !sLine.equals("---"))
                {
                sLine = queue.poll();
                }

            sLine = queue.poll();

            Map<String, LinkedHashMap> mapYaml = new LinkedHashMap<>();

            while (sLine != null)
                {
                String sName = UUID.randomUUID().toString();

                if (sLine.startsWith("# Source"))
                    {
                    sName = sLine.substring(sLine.lastIndexOf("/") + 1);
                    }

                StringBuilder sYaml = new StringBuilder();

                while(sLine != null && !sLine.equals("---") && !sLine.equals("(terminated)"))
                    {
                    if (!sLine.startsWith("# "))
                        {
                        sYaml.append(sLine).append('\n');
                        }

                    sLine = queue.poll();
                    }

                if (sYaml.length() > 0)
                    {
                    LinkedHashMap map = YAML_MAPPER.readValue(sYaml.toString(), LinkedHashMap.class);

                    mapYaml.put(sName, map);
                    }

                sLine = queue.poll();
                }

            return mapYaml;
            }
        catch (Throwable t)
            {
            throw new AssertionError("Failed to install chart", t);
            }
        }


    @SuppressWarnings("unchecked")
    private <T> T getYamlValue(Map map, String... asKey)
        {
        Object oResult = map;

        for (String sKey : asKey)
            {
            oResult = ((Map) oResult).get(sKey);
            }

        return (T) oResult;
        }

    private Map getNamedYamlValue(List list, String sName)
        {
        for (Object o : list)
            {
            if (o instanceof Map)
                {
                if (sName.equals(((Map) o).get("name")))
                    {
                    return (Map) o;
                    }
                }
            }

        return null;
        }


    // ----- data members ---------------------------------------------------

    @ClassRule
    public static TemporaryFolder s_tempFolder = new TemporaryFolder();

    private static File s_fileChart;

    private static final ObjectMapper YAML_MAPPER = new ObjectMapper(new YAMLFactory());

    /**
     * The yaml obtained from the initial defautl install to use to test default configurations.
     */
    private static Map<String, LinkedHashMap> s_mapDefaultYaml;
    }
