-- +migrate Up

INSERT INTO volunteer.volunteer(
	volunteer_id, display_name, contact_email, contact_phone, bio, joined)
	VALUES (1, 'gffg', 'flkgdf', '37489', 'fgkdshgfksdhf', '2020-06-10');

INSERT INTO volunteer.volunteer(
	volunteer_id, display_name, contact_email, contact_phone, bio, joined)
	VALUES (2, 'dfgd', 'flkgdf', '37489', 'fgkdshgfksdhf', '2020-06-10');

INSERT INTO volunteer.volunteer(
	volunteer_id, display_name, contact_email, contact_phone, bio, joined)
	VALUES (3, 'gfxvcxvxfg', 'flkgdf', '37489', 'fgkdshgfksdhf', '2020-06-10');

INSERT INTO organisation.organisation(
	organisation_id, display_name, locality, registration_no, contact_email, contact_phone, description, achievements, owner, type_of_org)
	VALUES (1, 'sfdf', 'fjdkf', 'dfjhj', 'hfjksd', 'sfdf', 'fjdkf', 'dfjhj', 1, 0);

INSERT INTO public.event_category(
	generated_id, name)
	VALUES (1, 'category 1');

INSERT INTO public.events(
	event_id, name, description, category)
	VALUES (1, 'test event', 'dfsjkdfh', 1);

INSERT INTO public.event_user_bridge(
	generated_id, event_id, volunteer_id, organisation_id)
	VALUES (0, 1, 1, 1);

INSERT INTO public.event_user_bridge(
	generated_id, event_id, volunteer_id, organisation_id)
	VALUES (1, 1, 2, 1);

INSERT INTO public.event_user_bridge(
	generated_id, event_id, volunteer_id, organisation_id)
	VALUES (2, 1, 3, 1);