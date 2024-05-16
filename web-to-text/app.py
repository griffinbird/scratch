import requests
import sys
import json
from bs4 import BeautifulSoup

# read a webpage and save the text as a string

webpage_1 = "https://learn.microsoft.com/en-us/azure/cosmos-db/mongodb/feature-support-50"
webpage_2 = "https://learn.microsoft.com/en-us/azure/cosmos-db/mongodb/feature-support-60"

def read_webpage(url):
    try:
        response = requests.get(url)
        response.raise_for_status()
        return response.text
    except requests.exceptions.RequestException as e:
        print(e)
        sys.exit(1)

webPageContents_1 = read_webpage(webpage_1)
webPageContents_2 = read_webpage(webpage_2)

html_doc = webPageContents_1

def extract_data(html_doc):
    response = requests.get(html_doc)
    soup = BeautifulSoup(response.text, 'html.parser')

    # Find the h3 tag and the next table
    h3_tag = soup.find('h3')
    table = h3_tag.find_next('table')

    # Extract the text from the h3 tag
    h3_text = h3_tag.get_text(strip=True)

    # Extract the table data
    table_data = []
    rows = table.find_all('tr')
    for row in rows:
        cells = row.find_all('td')
        cell_text = [td.get_text(strip=True) for td in cells]
        if len(cell_text) == 2:
            row_data = {'command': cell_text[0], 'supported': cell_text[1]}
            table_data.append(row_data)

    # Format the data into a dictionary
    data = {'title': h3_text, 'table_data': table_data}

    # Convert the dictionary to a JSON string
    json_data = json.dumps(data, indent=4)

    return json_data

webpage_1 = "https://learn.microsoft.com/en-us/azure/cosmos-db/mongodb/feature-support-50"
webpage_2 = "https://learn.microsoft.com/en-us/azure/cosmos-db/mongodb/feature-support-60"

json_data_1 = extract_data(webpage_1)
json_data_2 = extract_data(webpage_2)

print(json_data_1)
print(json_data_2)
