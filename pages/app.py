from flask import Flask 
import json
import requests
import os

app = Flask(__name__)

class Book:
    name = ""
    author = ""

    def __init__(self, name, author):
        self.name = name 
        self.author = author 

    def to_html_row(self):
        return "<tr><td>%s</td><td>%s</td></tr>" % (self.name, self.author)


def books_from_json(json_text):
    books_json = json.loads(json_text)
    result = []

    for book_info in books_json:
        book_name = book_info["bookName"]
        book_author = book_info["authorName"]

        result.append(Book(book_name, book_author))
    
    return result


@app.route("/")
def index():
    req = requests.get("http://books:80/api/books")
    books = books_from_json(req.text)
    rows = [book.to_html_row() for book in books]
    repl = "".join(rows)

    with open('./index.html') as html_file:
        html_text = html_file.read()
    
    html_text = html_text.replace("$BOOKS_INFO$", repl)
    return html_text

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=80)