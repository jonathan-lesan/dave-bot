import requests


def get_card(user_input):
    input_split = user_input.split(" ")[1:]
    input_split = ' '.join(input_split)
    set_name = ""

    if "(" in input_split and ")" in input_split:
        if len(input_split) > input_split.find(")")+1:
            return "Too many arguments. Please use the format: !card <card_name> (<set_name>)"
        card_name = input_split[0:input_split.find("(")]
        set_name = input_split[input_split.find("(")+1:input_split.find(")")]
        print(set_name)
    else:
        card_name = input_split

    r = requests.get('https://api.scryfall.com/cards/named?fuzzy=' + card_name)
    return_response = r.json()

    if return_response["object"] == "error":
        return return_response["details"]
    
    if set_name:
        return set_dive(set_name, return_response["prints_search_uri"])
    
    if return_response["layout"] == "transform" or return_response["layout"] == "modal_dfc":
        card_string = ""
        for face in return_response["card_faces"]:
            card_string += face["image_uris"]["large"] + " "
        return card_string
    else:
        return return_response["image_uris"]["large"]
        

def set_dive(set_name, api_uri):
    images_to_return = ""
    
    while True:
        r = requests.get(api_uri)
        request_response = r.json()
        matches = [x for x in request_response["data"] if x["set"] == set_name.lower() or set_name == x["set_name"]]

        if matches:
            for card in matches:
                if card["layout"] == "transform" or card["layout"] == "modal_dfc":
                    card_string = ""
                    for face in card["card_faces"]:
                        card_string += face["image_uris"]["large"] + " "
                    images_to_return += card_string
                else:
                    images_to_return += card["image_uris"]["large"] + " "
        
        if not request_response["has_more"]:
            break
        else:
            api_uri = request_response["next_page"]

    if not images_to_return:
        return "Card does not exist for set " + set_name
    else:
        return images_to_return


def set_lookup(set_name):
    r = requests.get('https://api.scryfall.com/sets')
    request_response = r.json()
    matches = [x for x in request_response["data"] if x["parent_set_code"] == set_name.lower()]
    if matches:
        return "yes"