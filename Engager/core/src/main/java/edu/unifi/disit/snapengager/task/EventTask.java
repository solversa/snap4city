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
package edu.unifi.disit.snapengager.task;

import java.io.IOException;
import java.util.Date;
import java.util.List;
import java.util.Locale;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.kie.api.builder.KieScanner;
import org.kie.api.runtime.KieContainer;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.scheduling.Trigger;
import org.springframework.scheduling.TriggerContext;
import org.springframework.scheduling.annotation.EnableScheduling;
import org.springframework.scheduling.annotation.SchedulingConfigurer;
import org.springframework.scheduling.config.ScheduledTaskRegistrar;
import org.springframework.scheduling.support.CronTrigger;
import org.springframework.stereotype.Component;

import edu.unifi.disit.snap4city.engager_utils.OrganizationType;
import edu.unifi.disit.snapengager.datamodel.profiledb.Event;
import edu.unifi.disit.snapengager.datamodel.profiledb.EventDAO;
import edu.unifi.disit.snapengager.datamodel.profiledb.SensorDAO;
import edu.unifi.disit.snapengager.exception.CredentialsException;
import edu.unifi.disit.snapengager.exception.UserprofileException;
import edu.unifi.disit.snapengager.service.IDataManagerService;

@EnableScheduling
@Component
public class EventTask implements SchedulingConfigurer {

	private static final Logger logger = LogManager.getLogger();

	@Value("${event.task.cron}")
	private String cronrefresh;

	KieContainer kContainer;
	KieScanner kScanner;

	@Autowired
	IDataManagerService dataService;

	@Autowired
	EventDAO eventRepo;

	@Autowired
	SensorDAO sensorRepo;

	@Override
	public void configureTasks(ScheduledTaskRegistrar taskRegistrar) {
		taskRegistrar.addTriggerTask(new Runnable() {
			@Override
			public void run() {
				try {
					logger.debug("-----------------------------------------------start execution for EVENT-----------------------------------------------");
					myTask();
				} catch (Exception e) {
					logger.error("Error catched {}", e);
				}
			}
		}, new Trigger() {
			@Override
			public Date nextExecutionTime(TriggerContext triggerContext) {
				CronTrigger trigger = new CronTrigger(cronrefresh);
				Date nextExec = trigger.nextExecutionTime(triggerContext);
				logger.debug("-----------------------------------------------next execution for EVENT will be fired at:{}-----------------------------------------------", nextExec);
				return nextExec;
			}
		});
	}

	private void myTask() throws UserprofileException, CredentialsException, IOException {
		Locale lang = new Locale("en");

		// delete events
		eventRepo.deleteAll();

		// popolate event from helsinki
		List<Event> helsinkiEvents = dataService.getEventData(OrganizationType.HELSINKI, lang);
		eventRepo.save(helsinkiEvents);

		// popolate event from DISIT/firenze
		List<Event> firenzeEvents = dataService.getEventData(OrganizationType.FIRENZE, lang);
		eventRepo.save(firenzeEvents);
	}
}