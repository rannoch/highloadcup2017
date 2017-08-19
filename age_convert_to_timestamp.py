from datetime import datetime
from dateutil.relativedelta import relativedelta
import calendar

now = datetime.now() - relativedelta(years=20)
timestamp = calendar.timegm(now.timetuple())

print timestamp