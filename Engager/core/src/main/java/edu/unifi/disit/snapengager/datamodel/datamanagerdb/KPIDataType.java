/* Snap4City Engager (SE)
   Copyright (C) 2015 DISIT Lab http://www.disit.org - University of Florence
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
package edu.unifi.disit.snapengager.datamodel.datamanagerdb;

public enum KPIDataType {
	S4CHelsinkiTrackerLocation("S4CHelsinkiTrackerLocation"), S4CAntwerpTrackerLocation("S4CAntwerpTrackerLocation"), S4CTuscanyTrackerLocation("S4CTuscanyTrackerLocation"), S4CHelsinkiAppUsage("S4CHelsinkiAppUsage"), S4CAntwerpAppUsage(
			"S4CAntwerpAppUsage"), S4CTuscanyAppUsage("S4CTuscanyAppUsage");

	private final String text;

	private KPIDataType(final String text) {
		this.text = text;
	}

	@Override
	public String toString() {
		return text;
	}
}