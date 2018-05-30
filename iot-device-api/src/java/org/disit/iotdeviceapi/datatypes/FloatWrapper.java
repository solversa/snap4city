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
public class FloatWrapper extends DataType {
    
    Float value;
    
    public FloatWrapper() {
        super();
        value = null;
    }
    
    @Override
    public String toString() {
        try  {
            return String.valueOf(value);
        }
        catch(Exception e) {
            status = Const.ERROR;
            return null;
        }
    }

    @Override
    public FloatWrapper fromString(String str) {
        try {
            FloatWrapper myFloat = new FloatWrapper();
            myFloat.setValue(Float.parseFloat(str));
            return myFloat;
        }
        catch(Exception e) {
            status = Const.ERROR;
            return null;
        }
    }

    public Float getValue() {
        return value;
    }

    public void setValue(Float value) {
        this.value = value;
    }
    
}
