(:require '[clojure.string :as str])

(defn wrap-to-json [list-name words-list]
  (json/generate-string {:metadata {:name list-name
                                    :size (count words-list)
                                    :packagedAt (java.util.Date/from (java.time.Instant/now))
                                    :version 1}

                         :words words-list} {:pretty true}))

(def exclusions '("“" "”" "gutenberg" "Gutenberg" "Petersburgh" "Mr." "Mrs." "—" "Chapter"))

(defn contains-exclusion? [str]
  (not-every? false? (map #(str/includes? str %) exclusions)))

(defn run [[book-uri list-name output-file]]
  (->> (slurp book-uri)
       (re-seq #"(?s)(.*?(?:\.|\?|!))(?: |$)")
       (map first)
       (filter #(and (> 100 (count %)) (< 15 (count %))))
       (filter  #(not (contains-exclusion? %)))
       (shuffle)
       (map #(str/replace % #"\s+" " "))
       (map #(str/replace % #"’" "'"))
       (map #(str/replace % #"_" ""))
       (map str/trim)
       (wrap-to-json list-name)
       (spit output-file)))

(run *command-line-args*)

(println "Done!")
