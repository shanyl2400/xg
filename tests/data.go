package tests

import "xg/entity"

var (
	students = []*entity.CreateStudentRequest{
		{
			Name:          GetFullName(),
			Gender:        false,
			Telephone:     RandTelephone(),
			Address:       RandArray(addresses),
			Email:         RandEmail(8),
			IntentSubject: RandArrayList(subjects, 3),
			Note:          RandString(12),
			OrderSourceID: 1,
		},
		{
			Name:          GetFullName(),
			Gender:        true,
			Telephone:     RandTelephone(),
			Address:       RandArray(addresses),
			Email:         RandEmail(8),
			IntentSubject: RandArrayList(subjects, 3),
			Note:          RandString(12),
			OrderSourceID: 1,
		},
		{
			Name:          GetFullName(),
			Gender:        false,
			Telephone:     RandTelephone(),
			Address:       RandArray(addresses),
			Email:         RandEmail(8),
			IntentSubject: RandArrayList(subjects, 3),
			Note:          RandString(12),
			OrderSourceID: 1,
		},
		{
			Name:          GetFullName(),
			Gender:        true,
			Telephone:     RandTelephone(),
			Address:       RandArray(addresses),
			Email:         RandEmail(8),
			IntentSubject: RandArrayList(subjects, 3),
			Note:          RandString(12),
			OrderSourceID: 1,
		},
		{
			Name:          GetFullName(),
			Gender:        false,
			Telephone:     RandTelephone(),
			Address:       RandArray(addresses),
			Email:         RandEmail(8),
			IntentSubject: RandArrayList(subjects, 3),
			Note:          RandString(12),
			OrderSourceID: 1,
		},
		{
			Name:          GetFullName(),
			Gender:        true,
			Telephone:     RandTelephone(),
			Address:       RandArray(addresses),
			Email:         RandEmail(8),
			IntentSubject: RandArrayList(subjects, 3),
			Note:          RandString(12),
			OrderSourceID: 1,
		},
		{
			Name:          GetFullName(),
			Gender:        false,
			Telephone:     RandTelephone(),
			Address:       RandArray(addresses),
			Email:         RandEmail(8),
			IntentSubject: RandArrayList(subjects, 3),
			Note:          RandString(12),
			OrderSourceID: 1,
		},
		{
			Name:          GetFullName(),
			Gender:        true,
			Telephone:     RandTelephone(),
			Address:       RandArray(addresses),
			Email:         RandEmail(8),
			IntentSubject: RandArrayList(subjects, 3),
			Note:          RandString(12),
			OrderSourceID: 1,
		},
	}

	orgs = []*entity.CreateOrgWithSubOrgsRequest{
		{
			OrgData: entity.CreateOrgRequest{
				Name:      "测试机构" + RandString(2),
				Address:   RandArray(addresses),
				AddressExt: RandString(10),
				Telephone: RandTelephone(),
			},
			SubOrgs: []*entity.CreateOrgRequest{
				{
					Name:      "校区" + RandString(2),
					Subjects:  RandArrayList(subjects, 2),
					Address:   RandArray(addresses),
					AddressExt: RandString(10),
					Telephone: "",
				},
				{
					Name:      "校区" + RandString(2),
					Subjects:  RandArrayList(subjects, 2),
					Address:   RandArray(addresses),
					AddressExt: RandString(10),
					Telephone: "",
				},
				{
					Name:      "校区" + RandString(2),
					Subjects:  RandArrayList(subjects, 2),
					Address:   RandArray(addresses),
					AddressExt: RandString(10),
					Telephone: "",
				},
			},
		},

		{
			OrgData: entity.CreateOrgRequest{
				Name:      "测试机构" + RandString(2),
				Address:   RandArray(addresses),
				AddressExt: RandString(10),
				Telephone: RandTelephone(),
			},
			SubOrgs: []*entity.CreateOrgRequest{
				{
					Name:      "校区" + RandString(2),
					Subjects:  RandArrayList(subjects, 2),
					Address:   RandArray(addresses),
					AddressExt: RandString(10),
					Telephone: "",
				},
				{
					Name:      "校区" + RandString(2),
					Subjects:  RandArrayList(subjects, 2),
					Address:   RandArray(addresses),
					AddressExt: RandString(10),
					Telephone: "",
				},
				{
					Name:      "校区" + RandString(2),
					Subjects:  RandArrayList(subjects, 2),
					Address:   RandArray(addresses),
					AddressExt: RandString(10),
					Telephone: "",
				},
			},
		},
		{
			OrgData: entity.CreateOrgRequest{
				Name:      "测试机构" + RandString(1),
				Address:   RandArray(addresses),
				AddressExt: RandString(10),
				Telephone: RandTelephone(),
			},
			SubOrgs: []*entity.CreateOrgRequest{
				{
					Name:      "校区" + RandString(2),
					Subjects:  RandArrayList(subjects, 2),
					Address:   RandArray(addresses),
					AddressExt: RandString(10),
					Telephone: "",
				},
			},
		},
		{
			OrgData: entity.CreateOrgRequest{
				Name:      "测试机构" + RandString(2),
				Address:   RandArray(addresses),
				AddressExt: RandString(10),
				Telephone: RandTelephone(),
			},
			SubOrgs: []*entity.CreateOrgRequest{
				{
					Name:      "校区" + RandString(2),
					Subjects:  RandArrayList(subjects, 2),
					Address:   RandArray(addresses),
					AddressExt: RandString(10),
					Telephone: "",
				},
				{
					Name:      "校区" + RandString(2),
					Subjects:  RandArrayList(subjects, 2),
					Address:   RandArray(addresses),
					AddressExt: RandString(10),
					Telephone: "",
				},
			},
		},
	}
)
