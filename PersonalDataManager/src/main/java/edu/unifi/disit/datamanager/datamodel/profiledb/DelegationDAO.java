/* Data Manager (DM).
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
package edu.unifi.disit.datamanager.datamodel.profiledb;

import java.util.Date;
import java.util.List;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;
import org.springframework.transaction.annotation.Transactional;

@Repository
// all with contrains: AndDeleteTimeIsNull
public interface DelegationDAO extends JpaRepository<Delegation, Long>, DelegationDAOCustom {

	List<Delegation> findByElementIdAndDeleteTimeIsNull(String elementId);

	List<Delegation> findByElementIdAndElementTypeAndDeleteTimeIsNull(String elementId, String elementType);

	List<Delegation> findByElementIdAndDeleteTimeIsNullAndUsernameDelegatedNotLike(String elementId, String usernameDelegated);

	Delegation findByIdAndDeleteTimeIsNull(Long delegationId);

	List<Delegation> findByDeleteTimeBefore(Date date);// used for proper delete

	Page<Delegation> findByElementIdAndDeleteTimeIsNull(String elementId, Pageable pageable);

	Page<Delegation> findByElementIdAndDeleteTimeIsNullAndUsernameDelegatedNotLike(String elementId, String usernameDelegated, Pageable pageable);

	@Modifying
	@Transactional
	@Query("delete from Delegation a where a.deleteTime < ?1")
	void deleteByDeleteTimeBefore(Date time);

	Page<Delegation> findByElementIdAndDeleteTimeIsNullAndUsernameDelegatedNotLikeAndUsernameDelegatedLike(String elementId, String usernameDelegated, String searchKey, Pageable pageable);

	List<Delegation> findByElementIdAndDeleteTimeIsNullAndUsernameDelegatedNotLikeAndUsernameDelegatedLike(String elementId, String usernameDelegated, String searchKey);
}