/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s.testing.spring;

import com.oracle.coherence.common.base.Logger;

import com.tangosol.util.MapEvent;
import com.tangosol.util.MapListener;

import org.springframework.context.annotation.Bean;

public class MapListenerBean<K, V>
        implements MapListener<K, V> {

    @Override
    public void entryInserted(MapEvent<K, V> mapEvent) {
        Logger.info("MapListenerBean.entryInserted "  + mapEvent);
    }

    @Override
    public void entryUpdated(MapEvent<K, V> mapEvent) {
        Logger.info("MapListenerBean.entryUpdated "  + mapEvent);
    }

    @Override
    public void entryDeleted(MapEvent<K, V> mapEvent) {
        Logger.info("MapListenerBean.entryDeleted "  + mapEvent);
    }
}
