(use  '[clojure.string])

(defn make-local-file-name [base-name]
  (str (replace (->> base-name
                     lower-case) #" " "-") ".json"))

(defn sentences-storage-path [base-name]
  (str "storage/sentences/" (make-local-file-name base-name)))

(defn words-storage-path [base-name]
  (str "storage/words/" (make-local-file-name base-name)))

(defn run-books-words [[remote-url name]]
  (clojure.java.shell/sh "bb" "book-word-list.clj" remote-url name (words-storage-path name)))

(defn run-books-sentences [[remote-url name]]
  (clojure.java.shell/sh "bb" "book-sentence-list.clj" remote-url name (sentences-storage-path name)))

(def sources [["https://www.gutenberg.org/files/84/84-0.txt" "Frankenstein"]
              ["https://www.gutenberg.org/ebooks/174.txt.utf-8" "Dorian Gray"]
              ["https://www.gutenberg.org/files/1342/1342-0.txt" "Pride and Prejudice"]
              ["https://www.gutenberg.org/files/1661/1661-0.txt" "Sherlock Holmes"]
              ["https://www.gutenberg.org/cache/epub/345/pg345.txt" "Dracula"]
              ["https://www.gutenberg.org/files/1952/1952-0.txt" "The Yellow Wallpaper"]
              ["https://www.gutenberg.org/files/98/98-0.txt" "A Tale of Two Cities"]
              ["https://www.gutenberg.org/cache/epub/64317/pg64317.txt" "The Great Gatsby"]
              ["https://www.gutenberg.org/cache/epub/1184/pg1184.txt" "The Count of Monte Cristo"]
              ["https://www.gutenberg.org/cache/epub/120/pg120.txt" "Treasure Island"]
              ["https://www.gutenberg.org/cache/epub/514/pg514.txt" "Little Women"]
              ["https://www.gutenberg.org/files/16/16-0.txt" "Peter Pan"]])

(defn run []
  (doall (map #((run-books-words %)
                (run-books-sentences %)) sources))

  (println "All done!"))

(run)

