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
import java.util.Hashtable;
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

import edu.unifi.disit.snapengager.datamodel.Data;
import edu.unifi.disit.snapengager.datamodel.datamanagerdb.KPIData;
import edu.unifi.disit.snapengager.datamodel.datamanagerdb.KPIDataDAO;
import edu.unifi.disit.snapengager.datamodel.datamanagerdb.KPIValueDAO;
import edu.unifi.disit.snapengager.datamodel.profiledb.Userprofile;
import edu.unifi.disit.snapengager.exception.CredentialsException;
import edu.unifi.disit.snapengager.exception.UserprofileException;
import edu.unifi.disit.snapengager.service.IDataManagerService;
import edu.unifi.disit.snapengager.service.IUserprofileService;

@EnableScheduling
@Component
public class LastLoginTask implements SchedulingConfigurer {

	private static final Logger logger = LogManager.getLogger();

	@Value("${lastlogin.task.cron}")
	private String cronrefresh;

	KieContainer kContainer;
	KieScanner kScanner;

	@Autowired
	IUserprofileService upservice;

	@Autowired
	IDataManagerService dataservice;

	@Autowired
	private KPIDataDAO kpidatarepo;

	@Autowired
	KPIValueDAO kvrepo;

	@Override
	public void configureTasks(ScheduledTaskRegistrar taskRegistrar) {
		taskRegistrar.addTriggerTask(new Runnable() {
			@Override
			public void run() {
				try {
					logger.debug("-----------------------------------------------start execution for LAST_LOGIN-----------------------------------------------");
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
				logger.debug("-----------------------------------------------next execution for LAST_LOGIN will be fired at:{}-----------------------------------------------", nextExec);
				return nextExec;
			}
		});
	}

	private void myTask() throws CredentialsException, IOException, UserprofileException {
		Locale lang = new Locale("en");
		Hashtable<String, Boolean> assistanceEnabled = dataservice.getAssistanceEnabled(lang);

		for (Data lastlogin : dataservice.getLastLoginData(lang)) {
			// check if the user has enable the assistance
			if ((assistanceEnabled.get(lastlogin.getUsername()) != null) && (assistanceEnabled.get(lastlogin.getUsername()))) {
				Userprofile up = upservice.get(lastlogin.getUsername(), lang);
				if (up == null)
					up = new Userprofile(lastlogin.getUsername());
				up.setLastlogin(new Date(Long.parseLong(lastlogin.getVariableValue())));
				upservice.save(up, lang);
			}
		}

		// update first login
		for (Userprofile up : upservice.getAll(lang)) {
			logger.debug("Analyze {}", up.getUsername());
			if (up.getFirstlogin() == null) {
				List<KPIData> kpidatas = kpidatarepo.findByUsernameAndValueNameContainingAndDeleteTimeIsNull(up.getUsername(), "AppUsage");
				Date firstlogin = up.getLastlogin();
				for (KPIData kpidata : kpidatas) {
					Date tempo = kvrepo.findTopByKpiIdOrderById(kpidata.getId()).getDataTime();
					if (tempo != null && tempo.before(firstlogin))
						firstlogin = tempo;
				}
				if (firstlogin != null) {
					up.setFirstlogin(firstlogin);
					upservice.save(up, lang);
				}
			}
		}
	}
}