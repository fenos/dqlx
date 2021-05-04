package testdata

func TestSchema() string {
	return `
		<NRHP>: default .
		<actor.dubbing_performances>: [uid] .
		<actor.film>: [uid] @count .
		<addr.city>: default .
		<addr.country>: default .
		<addr.housename>: default .
		<addr.housenumber>: default .
		<addr.postcode>: default .
		<addr.state>: default .
		<addr.street>: default .
		<alt_name>: default .
		<amenity>: default .
		<apple_movietrailer_id>: default .
		<architect>: default .
		<area>: default .
		<art_direction_by>: [uid] .
		<art_director.films_art_directed>: [uid] .
		<artist>: default .
		<artist_name>: default .
		<artwork_type>: default .
		<atm>: default .
		<attraction>: default .
		<barbecue_grill>: default .
		<barrier>: default .
		<board_type>: default .
		<building.colour>: default .
		<building.levels.underground>: default .
		<building.levels>: default .
		<building.material>: default .
		<building.part>: default .
		<building>: default .
		<car>: default .
		<casting_director.films_casting_directed>: [uid] .
		<casting_director>: [uid] .
		<character.portrayed_in_films>: [uid] .
		<character.portrayed_in_films_dubbed>: [uid] .
		<cinematographer.film>: [uid] .
		<cinematography>: [uid] .
		<collection.films_in_collection>: [uid] .
		<collections>: [uid] .
		<company.films>: [uid] .
		<company_role_or_service.companies_performing_this_role_or_service>: [uid] .
		<contact.email>: default .
		<contact.fax>: default .
		<contact.phone>: default .
		<contact.website>: default .
		<content_rating.country>: [uid] .
		<content_rating.minimum_accompanied_age>: default .
		<content_rating.minimum_unaccompanied_age>: default .
		<content_rating.rating_system>: [uid] .
		<content_rating_system.jurisdiction>: [uid] .
		<content_rating_system.ratings>: [uid] .
		<costume_design_by>: [uid] .
		<costumer_designer.costume_design_for_film>: [uid] .
		<country>: [uid] @reverse .
		<created_by>: default .
		<crew_gig.crew_role>: [uid] .
		<crew_gig.crewmember>: [uid] .
		<crew_gig.film>: [uid] .
		<crewmember.films_crewed>: [uid] .
		<cut.film>: [uid] .
		<cut.note>: string @lang .
		<cut.release_region>: [uid] .
		<cut.runtime>: default .
		<cut.type_of_cut>: [uid] .
		<description.en>: default .
		<description>: default .
		<designation>: default .
		<director.film>: [uid] @count @reverse .
		<distribution_medium.films_distributed_in_this_medium>: [uid] .
		<distributor.films_distributed>: [uid] .
		<distributors>: [uid] .
		<dubbing_performances>: [uid] .
		<edited_by>: [uid] .
		<editor.film>: [uid] .
		<ele>: default .
		<email>: string @index(exact) @upsert .
		<entrance>: default .
		<estimated_budget>: [uid] .
		<executive_produced_by>: [uid] .
		<fandango_id>: default .
		<fax>: default .
		<featured_locations>: [uid] .
		<featured_song.featured_in_film>: [uid] .
		<featured_song.performed_by>: [uid] .
		<featured_song>: [uid] .
		<fee>: default .
		<festival.date_founded>: datetime .
		<festival.focus>: [uid] .
		<festival.individual_festivals>: [uid] .
		<festival.location>: [uid] .
		<festival.sponsoring_organization>: [uid] .
		<festival_event.festival>: [uid] .
		<festival_event.films>: [uid] .
		<festival_focus.festivals_with_this_focus>: [uid] .
		<festival_sponsor.festivals_sponsored>: [uid] .
		<festivals>: [uid] .
		<filming>: [uid] .
		<fixme>: default .
		<format.format>: [uid] .
		<format>: [uid] .
		<former_name>: default .
		<friend>: [uid] .
		<genre>: [uid] @count @reverse .
		<gnis.county_id>: default .
		<gnis.county_name>: default .
		<gnis.created>: default .
		<gnis.feature_id>: default .
		<gnis.import_uuid>: default .
		<gnis.reviewed>: default .
		<gnis.state_id>: default .
		<gross_revenue>: [uid] .
		<height>: default .
		<heritage.operator>: default .
		<heritage>: default .
		<highway>: default .
		<hiking>: default .
		<historic.name>: default .
		<historic>: default .
		<history>: default .
		<http://www.w3.org/2000/01/rdf-schema#domain>: [uid] .
		<http://www.w3.org/2000/01/rdf-schema#range>: [uid] .
		<http://www.w3.org/2002/07/owl#inverseOf>: [uid] .
		<id>: default .
		<information>: default .
		<initial_release_date>: datetime @index(year) .
		<internet_access.fee>: default .
		<internet_access>: default .
		<job.films_with_this_crew_job>: [uid] .
		<landuse>: default .
		<language>: [uid] .
		<layer>: default .
		<leisure>: default .
		<level>: default .
		<levels>: default .
		<loc>: geo @index(geo) .
		<location.featured_in_films>: [uid] .
		<locations>: [uid] .
		<man_made>: default .
		<map_size>: default .
		<map_type>: default .
		<metacritic_id>: default .
		<mobile>: default .
		<music>: [uid] .
		<music_contributor.film>: [uid] .
		<name.ca>: default .
		<name.da>: default .
		<name.en>: default .
		<name.es>: default .
		<name.fr>: default .
		<name.ko>: default .
		<name.ru>: default .
		<name.sv>: default .
		<name.zh>: default .
		<name>: string @index(fulltext, hash, term, trigram) @lang .
		<name_2>: default .
		<name_alt>: default .
		<natural>: default .
		<netflix_id>: default .
		<note>: default .
		<nytimes_id>: default .
		<office>: default .
		<old_name>: default .
		<opening_hours>: default .
		<operator>: default .
		<other_companies>: [uid] .
		<other_crew>: [uid] .
		<pass>: password .
		<payment.bitcoin>: default .
		<payment.credit_cards>: default .
		<performance.actor>: [uid] .
		<performance.character>: [uid] .
		<performance.character_note>: string @lang .
		<performance.film>: [uid] .
		<performance.special_performance_type>: [uid] .
		<person_or_entity_appearing_in_films>: [uid] .
		<personal_appearance.film>: [uid] .
		<personal_appearance.person>: [uid] .
		<personal_appearance.type_of_appearance>: [uid] .
		<personal_appearance_type.appearances>: [uid] .
		<personal_appearances>: [uid] .
		<phone>: default .
		<plain_cut.note>: string .
		<plain_name>: string @index(fulltext, hash, term, trigram) .
		<plain_performance.character_note>: string .
		<plain_tagline>: string .
		<post_production>: [uid] .
		<pre_production>: [uid] .
		<prequel>: [uid] .
		<primary_language>: [uid] .
		<produced_by>: [uid] .
		<producer.film>: [uid] .
		<producer.films_executive_produced>: [uid] .
		<production_companies>: [uid] .
		<production_company.films>: [uid] .
		<production_design_by>: [uid] .
		<production_designer.films_production_designed>: [uid] .
		<railway>: default .
		<rated>: [uid] @reverse .
		<rating>: [uid] @reverse .
		<ref.nrhp>: default .
		<ref>: default .
		<regional_release_date.release_date>: datetime .
		<regional_release_date.release_region>: [uid] .
		<release_date_s>: [uid] .
		<roof.height>: default .
		<roof.material>: default .
		<roof.shape>: default .
		<rooms>: default .
		<rottentomatoes_id>: default .
		<runtime>: [uid] .
		<sequel>: [uid] .
		<series.films_in_series>: [uid] .
		<series>: [uid] .
		<set_decoration_by>: [uid] .
		<set_designer.sets_designed>: [uid] .
		<sfgov.OBJNAME>: default .
		<ship.type>: default .
		<shop>: default .
		<smoking>: default .
		<song.films>: [uid] .
		<song_performer.songs>: [uid] .
		<songs>: [uid] .
		<soundtrack>: [uid] .
		<source.opening_hours>: default .
		<source.start_date>: default .
		<source>: default .
		<special_performance_type.performance_type>: [uid] .
		<species>: default .
		<starring>: [uid] @count .
		<stars>: default .
		<start_date>: default .
		<story_by>: [uid] .
		<story_contributor.story_credits>: [uid] .
		<subject.films>: [uid] .
		<subjects>: [uid] .
		<surface>: default .
		<tagline>: string @lang .
		<theatre>: default .
		<toilets.wheelchair>: default .
		<topic_server.schemastaging_corresponding_entities_type>: [uid] .
		<topic_server.webref_cluster_members_type>: [uid] .
		<tourism>: default .
		<tourism_2>: default .
		<traileraddict_id>: default .
		<trailers>: [uid] .
		<type.property.expected_type>: [uid] .
		<type.property.reverse_property>: [uid] .
		<type.property.schema>: [uid] .
		<url>: default .
		<verified>: default .
		<website>: default .
		<wheelchair>: default .
		<wifi>: default .
		<wikidata>: default .
		<wikipedia>: default .
		<wpt_description>: default .
		<writer.film>: [uid] .
		<written_by>: [uid] .
		<zoo>: default .
		type <Actor> {
			name
			actor.film
			actor.dubbing_performances
		}
		type <ArtDirector> {
			name
			art_director.films_art_directed
		}
		type <CastingDirector> {
			name
			casting_director.films_casting_directed
		}
		type <Character> {
			name
			character.portrayed_in_films
			character.portrayed_in_films_dubbed
		}
		type <Cinematographer> {
			name
			cinematographer.film
		}
		type <Collection> {
			name
			collection.films_in_collection
		}
		type <Company> {
			name
			company.films
		}
		type <CompanyRoleOrService> {
			name
			company_role_or_service.companies_performing_this_role_or_service
		}
		type <CostumeDesigner> {
			name
			costumer_designer.costume_design_for_film
		}
		type <CrewGig> {
			name
			crew_gig.crew_role
			crew_gig.crewmember
			crew_gig.film
		}
		type <CrewMember> {
			name
			crewmember.films_crewed
		}
		type <Critic> {
			name
		}
		type <Cut> {
			name
			cut.film
			cut.note
			cut.release_region
			cut.runtime
			cut.type_of_cut
		}
		type <CutType> {
			name
		}
		type <Director> {
			name
			director.film
		}
		type <DistributionMedium> {
			name
			distribution_medium.films_distributed_in_this_medium
		}
		type <Distributor> {
			name
			distributor.films_distributed
		}
		type <Editor> {
			name
			editor.film
		}
		type <FeaturedSong> {
			name
			featured_song.featured_in_film
			featured_song.performed_by
		}
		type <Festival> {
			name
			festival.date_founded
			festival.focus
			festival.individual_festivals
			festival.location
			festival.sponsoring_organization
		}
		type <FestivalEvent> {
			name
			festival_event.festival
			festival_event.films
		}
		type <FestivalFocus> {
			name
			festival_focus.festivals_with_this_focus
		}
		type <FestivalSponsor> {
			name
			festival_sponsor.festivals_sponsored
		}
		type <FestivalSponsorship> {
			name
			festival_sponsor.festivals_sponsored
		}
		type <Film> {
			apple_movietrailer_id
			art_direction_by
			casting_director
			cinematography
			collections
			costume_design_by
			country
			distributors
			dubbing_performances
			edited_by
			executive_produced_by
			fandango_id
			featured_locations
			featured_song
			festivals
			format
			genre
			initial_release_date
			locations
			metacritic_id
			music
			name
			netflix_id
			personal_appearances
			prequel
			produced_by
			production_companies
			production_design_by
			rating
			release_date_s
			rottentomatoes_id
			sequel
			series
			set_decoration_by
			songs
			starring
			story_by
			subjects
			tagline
			traileraddict_id
			written_by
			post_production
			pre_production
			runtime
			other_crew
			other_companies
			primary_language
			soundtrack
			trailers
			gross_revenue
			estimated_budget
			filming
			language
		}
		type <Format> {
			name
			format.format
		}
		type <Generic> {
			name
		}
		type <Genre> {
			name
		}
		type <Job> {
			name
			job.films_with_this_crew_job
		}
		type <Location> {
			name
			location.featured_in_films
		}
		type <MusicContributor> {
			name
			music_contributor.film
		}
		type <Performance> {
			performance.actor
			performance.character
			performance.character_note
			performance.film
			performance.special_performance_type
		}
		type <PersonOrEntityAppearingInFilm> {
			name
			person_or_entity_appearing_in_films
			personal_appearance.film
		}
		type <PersonalAppearance> {
			name
			personal_appearance.film
			personal_appearance.person
			personal_appearance.type_of_appearance
			personal_appearance_type.appearances
		}
		type <PersonalAppearanceType> {
			name
		}
		type <Producer> {
			name
			producer.film
			producer.films_executive_produced
		}
		type <ProductionCompany> {
			name
			production_company.films
		}
		type <ProductionDesigner> {
			name
			production_designer.films_production_designed
		}
		type <Rating> {
			name
			content_rating.country
			content_rating.minimum_accompanied_age
			content_rating.minimum_unaccompanied_age
			content_rating.rating_system
		}
		type <RatingSystem> {
			name
			content_rating_system.ratings
			content_rating_system.jurisdiction
		}
		type <RegionalReleaseDate> {
			name
			regional_release_date.release_date
			regional_release_date.release_region
		}
		type <Series> {
			name
			series.films_in_series
		}
		type <SetDecorator> {
			name
			set_designer.sets_designed
		}
		type <Song> {
			name
		}
		type <SongPerformer> {
			name
			song_performer.songs
		}
		type <SpecialPerformanceType> {
			name
			special_performance_type.performance_type
		}
		type <StoryContributor> {
			name
			story_contributor.story_credits
		}
		type <Subject> {
			name
			subject.films
		}
		type <Writer> {
			name
		}
`
}
