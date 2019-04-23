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
package edu.unifi.disit.snapengager;

import org.apache.http.client.config.RequestConfig;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClientBuilder;
import org.apache.http.impl.conn.PoolingHttpClientConnectionManager;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.builder.SpringApplicationBuilder;
import org.springframework.boot.web.support.SpringBootServletInitializer;
import org.springframework.context.annotation.Bean;
import org.springframework.http.client.ClientHttpRequestFactory;
import org.springframework.http.client.HttpComponentsClientHttpRequestFactory;
import org.springframework.ldap.core.LdapTemplate;
import org.springframework.ldap.core.support.LdapContextSource;
import org.springframework.scheduling.TaskScheduler;
import org.springframework.scheduling.concurrent.ConcurrentTaskScheduler;
import org.springframework.web.client.RestTemplate;

import edu.unifi.disit.snapengager.datamodel.ldap.LDAPUserDAO;
import edu.unifi.disit.snapengager.datamodel.ldap.LDAPUserDAOImpl;

@SpringBootApplication

public class Application extends SpringBootServletInitializer {

	@Value("${spring.ldap.url}")
	private String ldapUrl;

	@Value("${spring.ldap.basicdn}")
	private String ldapBasicDN;

	public static void main(String[] args) throws Throwable {
		SpringApplication.run(Application.class, args);
	}

	// to enable scenario with my external tomcat
	private static Class<Application> applicationClass = Application.class;

	@Override
	protected SpringApplicationBuilder configure(SpringApplicationBuilder application) {
		return application.sources(applicationClass);
	}

	// to avoid annoying DEBUG message about missing executor for Task (Could not find default ScheduledExecutorService bean)
	@Bean
	public TaskScheduler taskScheduler() {
		return new ConcurrentTaskScheduler(); // single threaded by default
	}

	@Bean
	public LdapContextSource contextSource() {
		LdapContextSource ctxSrc = new LdapContextSource();
		ctxSrc.setUrl(ldapUrl);
		ctxSrc.setBase(ldapBasicDN);
		ctxSrc.setAnonymousReadOnly(true);
		return ctxSrc;
	}

	@Bean
	public LdapTemplate ldapTemplate() {
		return new LdapTemplate(contextSource());
	}

	@Bean
	public LDAPUserDAO ldapUserDAO() {
		return new LDAPUserDAOImpl();
	}

	@Bean
	public ClientHttpRequestFactory createRequestFactory(@Value("${connection.timeout}") String timeout) {
		PoolingHttpClientConnectionManager connectionManager = new PoolingHttpClientConnectionManager();
		// connectionManager.setMaxTotal(maxTotalConn);
		// connectionManager.setDefaultMaxPerRoute(maxPerChannel);

		RequestConfig config = RequestConfig.custom().
		/* */setConnectTimeout(Integer.parseInt(timeout)).
		/* */setConnectionRequestTimeout(Integer.parseInt(timeout)).build();

		CloseableHttpClient httpClient = HttpClientBuilder.create().setConnectionManager(connectionManager).setDefaultRequestConfig(config).build();
		return new HttpComponentsClientHttpRequestFactory(httpClient);
	}

	@Bean
	public RestTemplate restTemplate(ClientHttpRequestFactory factory) {
		return new RestTemplate();
	}
}