import sys

from PyPDF2 import PdfFileReader, PdfFileWriter
from PyPDF2.generic import BooleanObject, NameObject, IndirectObject

pdf = sys.argv[1]
reader = PdfFileReader(pdf)

for key in reader.getFormTextFields().keys():
  print(key)

values = {
  'Last, First, Middle': 'Testerson, Testy, S',
  'Service Branch': 'Air Force',
  'Trusted Agent Name': 'Another Person',
    }

writer = PdfFileWriter()
if "/AcroForm" not in writer._root_object:
    writer._root_object.update({
        NameObject("/AcroForm"): IndirectObject(len(writer._objects), 0, writer)})
writer._root_object["/AcroForm"].update(
    {NameObject("/NeedAppearances"): BooleanObject(True)})

def update_fields(page):
  if '/Annots' in page:
    writer.updatePageFormFieldValues(page, values)

writer.appendPagesFromReader(reader, update_fields)

out = open(sys.argv[2], "wb")
writer.write(out)
out.close()
