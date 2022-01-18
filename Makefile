USER_PATH = './user_srv/userpb'
DETAILS_PATH = './user_details_srv/detailspb'

user:
	cd $(USER_PATH) && \
		buf generate
