# 冒泡算法
'''
def bubble_sort(alist):
    n = len(alist)
    for j in range(n - 1):
        count = 0
        for i in range(n - 1 - j):
            if alist[i] > alist[i + 1]:
                (alist[i], alist[i + 1]) = (alist[i + 1], alist[i])
                count += 1
        if count == 0:
            return
'''

# 先序排列算法
'''
def preorder(self, node):
    if node is None:
        return
    print(node.elem, end=' ')
    self.preorder(node.lchild)
    self.preorder(node.rchild)
'''

# 获取一个网页
'''
def get_one_page(url):
    try:
        headers = {'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36'}
        response = requests.get(url, headers=headers)
        if response.status_code == 200:
            return response.text
        return None
    except RequestException:
        return None
'''

# import datetime
# def parse_expenses(expenses_string):
#     """Parse the list of expenses and return the list of triples (date, value, currency).
#     Ignore lines starting with #.
#     Parse the date using datetime.
#     Example expenses_string:
#         2016-01-02 -34.01 USD
#         2016-01-03 2.59 DKK
#         2016-01-03 -2.72 EUR
#     """
'''
    expenses = []
    for line in expenses_string.split('\n'):
        if line.startswith('#'):
            continue
        (date, value, currency) = line.split(' -')
        date = datetime.strptime(date, '%Y-%m-%d')
        value = float(value)
        currency = currency.strip()
        expenses.append((date, value, currency))
    return expenses
'''