package sheetsmock

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
)

func (r Form) CreateCSV(log *zap.SugaredLogger) {
	name := fmt.Sprintf("survey-%d.%s", r.ID, "csv")
	// Create a new CSV file
	file, err := os.Create(name)
	if err != nil {
		log.Errorln("Error creating CSV file:", err)
		return
	}
	defer file.Close()

	log.Infoln("CSV file created successfully.", name)
}

func (r Question) UpdateFirstRow(log *zap.SugaredLogger) {
	name := fmt.Sprintf("survey-%d.%s", r.FormID, "csv")
	// Read CSV file
	records := readCSV(log, name)

	if len(records) == 0 {
		newRow := getNewRowQuestion(r)
		records = append(records, newRow)
	}

	// Write updated records back to CSV file
	writeCSV(log, name, records)

	log.Infoln("CSV file updated successfully.", name)
}

func writeCSV(log *zap.SugaredLogger, filename string, records [][]string) {
	// Create a new CSV file
	file, err := os.Create(filename)
	if err != nil {
		log.Errorln("Error creating CSV file:", err)
		return
	}
	defer file.Close()

	// Create a new CSV writer
	writer := csv.NewWriter(file)

	// Write rows to the CSV file
	err = writer.WriteAll(records)
	if err != nil {
		log.Errorln("Error writing CSV file:", err)
		return
	}

	// Flush the writer to write the data to the file
	writer.Flush()

	if err := writer.Error(); err != nil {
		log.Infoln("Error flushing CSV writer:", err)
		return
	}

	log.Infoln("CSV file updated successfully.")
}

func readCSV(log *zap.SugaredLogger, filename string) [][]string {
	var file *os.File
	var err error
	// Open the CSV file
	file, err = os.Open(filename)
	if err != nil {
		log.Errorln("Error opening CSV file:", err)
		createCSV(log, filename)
		file, err = os.Open(filename)
		if err != nil {
			log.Errorln(err)
		}
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all records from the CSV file
	records, err := reader.ReadAll()
	if err != nil {
		log.Errorln("Error reading CSV file:", err)
		return nil
	}

	return records
}

func createCSV(log *zap.SugaredLogger, name string) {
	// Create a new CSV file
	file, err := os.Create(name)
	if err != nil {
		log.Errorln("Error creating CSV file:", err)
		return
	}
	defer file.Close()

	log.Infoln("CSV file created successfully.", name)
}

func (r Response) UpdateCSV(log *zap.SugaredLogger) {
	name := fmt.Sprintf("survey-%d.%s", r.FormID, "csv")
	// Open the CSV file in append mode
	file, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Errorln("Error opening CSV file:", err)
		return
	}
	defer file.Close()

	// Create a new CSV writer
	writer := csv.NewWriter(file)

	// Write a new row to the CSV file
	newRow := []string{fmt.Sprintf("%d", r.ResponseID), r.RespondentID, r.CreatedAt.GoString(), r.UpdatedAt.GoString()}
	writer.Write(newRow)

	// Flush the writer to write the data to the file
	writer.Flush()

	if err := writer.Error(); err != nil {
		log.Infoln("Error updating CSV file:", err)
		return
	}

	log.Infoln("CSV file updated successfully.")
}

func (r AnswerEvent) UpdateCSV(log *zap.SugaredLogger) {
	name := fmt.Sprintf("survey-%d.%s", r.Response.FormID, "csv")
	records := readCSV(log, name)

	var found bool
	var row int
	var anscol int = -1
	for i, record := range records {
		if i == 0 {
			continue
		}
		if record[0] == fmt.Sprintf("%d", r.Response.ResponseID) {
			found = true
			row = i
			for j, col := range records[0] {
				var queID string
				ques := strings.Split(col, ":{")
				if len(ques) > 0 {
					queID = ques[0]
				}
				if queID == fmt.Sprintf("%d", r.Answer.QuestionID) {
					anscol = j
					break
				}
			}
			break
		}
	}
	// create new row
	if !found {
		// Read CSV file
		file, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			log.Errorln("Error opening CSV file:", err)
			return
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		newRow := getNewAnswerRow(r.Response, r.Answer, records[0])
		// Write a new row to the CSV file
		writer.Write(newRow)
		// Flush the writer to write the data to the file
		writer.Flush()
		if err := writer.Error(); err != nil {
			log.Infoln("Error updating CSV file:", err)
			return
		}
		log.Infoln("CSV file row added successfully.")
	} else {
		// update row
		records[row][3] = r.Answer.UpdatedAt.GoString()
		if r.Answer.AnswerText != nil {
			records[row][anscol] = *r.Answer.AnswerText
		} else {
			records[row][anscol] = fmt.Sprintf("%d", *r.Answer.AnswerOptionID)
		}
		writeCSV(log, name, records)
		log.Infoln("CSV file row updated successfully.")
	}

}

func getNewAnswerRow(res Response, ans Answer, columns []string) []string {
	newRow := make([]string, len(columns))
	newRow[0], newRow[1], newRow[2], newRow[3] = fmt.Sprintf("%d", res.ResponseID), res.RespondentID, res.CreatedAt.GoString(), ans.UpdatedAt.GoString()
	for i, v := range columns {
		var queID string
		ques := strings.Split(v, ":{")
		if len(ques) > 0 {
			queID = ques[0]
		}
		if queID == fmt.Sprintf("%d", ans.QuestionID) {
			if ans.AnswerText != nil {
				newRow[i] = *ans.AnswerText
			} else {
				newRow[i] = fmt.Sprintf("%d", *ans.AnswerOptionID)
			}
			break
		}
	}
	return newRow
}

func getNewRowQuestion(r Question) []string {
	newRow := make([]string, 5)
	newRow[0], newRow[1], newRow[2], newRow[3] = "ResponseID", "RespondentID", "CreatedAt", "UpdatedAt"
	newRow[4] = r.Column()
	return newRow
}
