/* Data Manager (DM).
   Copyright (C) 2015 DISIT Lab http://www.disit.org - University of Florence
   This program is free software; you can redistribute it and/or
   modify it under the terms of the GNU General Public License
   as published by the Free Software Foundation; either version 2
   of the License, or (at your option) any later version.
   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.
   You should have received a copy of the GNU General Public License
   along with this program; if not, write to the Free Software
   Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301, USA. */
package edu.unifi.disit.datamanager.datamodel.profiledb;

import java.io.IOException;

import com.fasterxml.jackson.core.JsonGenerator;
import com.fasterxml.jackson.databind.SerializerProvider;
import com.fasterxml.jackson.databind.ser.std.StdSerializer;

public class KPIMetadataSerializer extends StdSerializer<KPIMetadata> {

	private static final long serialVersionUID = 1L;

	public KPIMetadataSerializer() {
		this(null);
	}

	public KPIMetadataSerializer(Class<KPIMetadata> t) {
		super(t);
	}

	@Override
	public void serialize(KPIMetadata kpiMetadata, JsonGenerator jgen, SerializerProvider provider) throws IOException {

		jgen.writeStartObject();

		if (kpiMetadata.getId() != null) {
			jgen.writeNumberField("id", kpiMetadata.getId());
		}
		if (kpiMetadata.getKpiId() != null) {
			jgen.writeNumberField("kpiId", kpiMetadata.getKpiId());
		}
		if (kpiMetadata.getDeleteTime() != null) {
			jgen.writeNumberField("deleteTime", kpiMetadata.getDeleteTime().getTime());
		}
		if (kpiMetadata.getElapseTime() != null) {
			jgen.writeNumberField("elapseTime", kpiMetadata.getElapseTime().getTime());
		}
		if (kpiMetadata.getKey() != null) {
			jgen.writeStringField("key", kpiMetadata.getKey());
		}
		if (kpiMetadata.getValue() != null) {
			jgen.writeStringField("value", kpiMetadata.getValue());
		}

		jgen.writeEndObject();
	}
}