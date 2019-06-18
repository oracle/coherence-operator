package com.oracle.coherence.examples;

import com.tangosol.net.events.EventInterceptor;
import com.tangosol.net.events.annotation.Interceptor;
import com.tangosol.net.events.partition.cache.CacheLifecycleEvent;

import com.tangosol.coherence.metrics.internal.MetricSet;

import java.io.IOException;
import java.io.OutputStream;
import java.io.Serializable;
import java.net.InetSocketAddress;
import java.util.Date;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.Map;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import com.sun.net.httpserver.HttpServer;

/**
 * Test EventInterceptor to allow prometheus to scrape.
 */
@Interceptor(identifier = "Custom Prometheus", cacheLifecycleEvents = {CacheLifecycleEvent.Type.CREATED, CacheLifecycleEvent.Type.DESTROYED})
public class PromInterceptor
        implements EventInterceptor<CacheLifecycleEvent>, Serializable
    {

    // ----- constructors ---------------------------------------------------

    /**
     * Construct a MutatingInterceptor that will register for all mutable events.
     */
    public PromInterceptor()
        {
        super();

        System.out.println("PromInterceptor: calling constructor...");
        f_mapMetricsValueNameToHelp.put("total", "Total Number of Coherence JMX Mbean Scrapes");
        f_mapMetricsValueNameToHelp.put("duration_ms", "Duration of Coherence JMX Mbean Scrapes in milliseconds");
        }

    // ----- EventInterceptor methods ---------------------------------------

    /**
     * {@inheritDoc}
     */
    @Override
    public void onEvent(CacheLifecycleEvent event)
        {
        String sCache = event.getCacheName();
        switch (event.getType())
            {
            case CREATED:
            try
                {
                System.out.println("Creating a custom metrics endpoint for cache " + sCache + "...");
                m_server = HttpServer.create(new InetSocketAddress(8200), 0);
                m_server.createContext("/metrics", new MyHandler());
                m_server.setExecutor(null); // creates a default executor
                m_server.start();
                System.out.println("Custom metrics HTTP server started.");
                }
            catch(IOException ioe)
                {
                System.out.println("PromInterceptor.onEvent: http server start failed with exception: " + ioe.getMessage());
                }
                break;

            case DESTROYED:
                System.out.println("PromInterceptor.onEvent: stopping http server...");
                m_server.stop(1);
                break;
            }
        }

    static class MyHandler implements HttpHandler
        {
        @Override
        public void handle(HttpExchange t) throws IOException
            {
            System.out.println("PromInterceptor: http server handler handle request");
            Map<String, Number> mapMetricsNameToValue = new HashMap<>();
            mapMetricsNameToValue.put("total PromInterceptor updates", 100);
            Map<String, String> mapTags = new HashMap<>();
            mapTags.put("PromInterceptor", new Date().toString());
            MetricSet collector = new MetricSet("PromInterceptor_custom_scrape", 0,
                    mapMetricsNameToValue, mapTags, PromInterceptor.f_mapMetricsValueNameToHelp);

            String response = MetricSet.toMetricLines(collector);
            System.out.println("PromInterceptor: http server handler response: " + response);
            t.sendResponseHeaders(200, response.length());
            OutputStream os = t.getResponseBody();
            os.write(response.getBytes());
            os.close();
            System.out.println("PromInterceptor: http server handler handle end");
            }
        }

    /**
     * Map coherence_jmx_scrape metric value name to its help description.
     */
    protected static final Map<String, String> f_mapMetricsValueNameToHelp = new LinkedHashMap<>();

    private HttpServer m_server;
    }
