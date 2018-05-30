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
package org.disit.iotdeviceapi.datatypes;

import org.disit.iotdeviceapi.utils.Const;

/**
 * 
 * @author Mirco Soderi @ DISIT DINFO UNIFI (mirco.soderi at unifi dot it)
 */
public class StringWrapper extends DataType {
    
    String value;
    
    public StringWrapper() {
        super();
        value = null;
    }
    
    @Override
    public String toString() {
        try {
            return value;
        } 
        catch(Exception e) {
            status = Const.ERROR;
            return null;
        }
    }
    
    @Override
    public StringWrapper fromString(String str) {
        try {
            StringWrapper myStr = new StringWrapper();
            myStr.setValue(str);
            return myStr;
        }
        catch(Exception e) {
            status = Const.ERROR;
            return null;
        }
    }

    public String getValue() {
        return value;
    }

    public void setValue(String value) {
        this.value = value;
    }

}
