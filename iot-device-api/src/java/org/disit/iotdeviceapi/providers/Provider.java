/* IOTDEVICEAPI
   Copyright (C) 2017 DISIT Lab http://www.disit.org - University of Florence

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as
   published by the Free Software Foundation, either version 3 of the
   License, or (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>. */
package org.disit.iotdeviceapi.providers;

import org.disit.iotdeviceapi.repos.Repos;
import org.disit.iotdeviceapi.utils.Const;
import org.disit.iotdeviceapi.logging.XLogger;

/**
 * 
 * @author Mirco Soderi @ DISIT DINFO UNIFI (mirco.soderi at unifi dot it)
 */
public class Provider {
    
    private final Repos datasource;
    private int status;
    private XLogger xlogger;
    
    public Provider(Repos datasource) {
        this.datasource = datasource;
        this.status = Const.OK;
    }

    public Repos getDatasource() {
        return datasource;
    }
    
    public Object[] get(Object[] queries) {
        return null;
    }

    public int getStatus() {
        return status;
    }

    public void setStatus(int status) {
        this.status = status;
    }

    public XLogger getXlogger() {
        return xlogger;
    }

    public void setXlogger(XLogger xlogger) {
        this.xlogger = xlogger;
    }
    
}
