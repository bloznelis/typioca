(defn wrap-to-json [words-list]
  (json/generate-string {:metadata {:name "common-words"
                                    :size (count words-list)
                                    :packagedAt (java.util.Date/from (java.time.Instant/now))
                                    :version 1}

                         :words words-list} {:pretty true}))

;; Expects text file which words are separated by new lines
(defn run [[list-uri output-file]]
  (->>
   (clojure.string/split (slurp list-uri) #"\n")
   (map clojure.string/lower-case)
   (filter #(and (> 7 (count %)) (< 1 (count %))))
   (wrap-to-json)
   (spit output-file)))

(run *command-line-args*)

(println "Done!")
