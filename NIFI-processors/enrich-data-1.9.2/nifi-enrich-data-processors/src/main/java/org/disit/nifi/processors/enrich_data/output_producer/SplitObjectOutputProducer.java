/**
 *  Nifi EnrichData processor
 *  
 *  Copyright (C) 2020 DISIT Lab http://www.disit.org - University of Florence
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU Affero General Public License as
 *  published by the Free Software Foundation, either version 3 of the
 *  License, or (at your option) any later version.
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU Affero General Public License for more details.
 *  You should have received a copy of the GNU Affero General Public License
 *  along with this program.  If not, see <http://www.gnu.org/licenses/>. 
 */


package org.disit.nifi.processors.enrich_data.output_producer;

import java.io.IOException;
import java.io.OutputStream;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

import org.apache.nifi.flowfile.FlowFile;
import org.apache.nifi.processor.ProcessSession;
import org.apache.nifi.processor.io.OutputStreamCallback;

import com.google.gson.JsonElement;
import com.google.gson.JsonObject;
/**
 * This OutputProducer implementation produces 
 * a flow file per member. 
 *
 */
public class SplitObjectOutputProducer implements OutputProducer {

	@Override
	public List<FlowFile> produceOutput(JsonObject rootObj, FlowFile inFlowFile, final ProcessSession session) {		
		List<FlowFile> outputList = new ArrayList<>();
		
		rootObj.entrySet().stream().forEach( (Map.Entry<String , JsonElement> rootEntry) -> { 
			
			FlowFile ff = session.create( inFlowFile );
			ff = session.write( ff , new OutputStreamCallback() {
				
				@Override
				public void process(OutputStream out) throws IOException {
					out.write( rootEntry.getValue().getAsJsonObject().toString().getBytes() );
				}
			});
			
			outputList.add( ff );
			
		});
		
		session.remove( inFlowFile );
		
		return outputList;
	}

}
